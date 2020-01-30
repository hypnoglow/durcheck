package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLintPackage(t *testing.T) {
	// buffer to send output to
	buf := &bytes.Buffer{}

	// exit code to capture
	exitCode := -1
	exitFunc := func(code int) {
		exitCode = code
	}

	stderr = buf
	exit = exitFunc
	os.Args = []string{"", "./testdata/foo"}

	main()

	expectedOutput := `testdata/foo/foo.go:7:2: implicit time.Duration means nanoseconds in "time.Sleep(60)"`
	actualOutput := strings.TrimSpace(buf.String())
	if actualOutput != expectedOutput {
		t.Errorf("Expected output to be %q but got %q", expectedOutput, actualOutput)
	}

	expectedCode := 1
	if exitCode != expectedCode {
		t.Errorf("Expected exit code to be %v but got %v", expectedCode, exitCode)
	}
}
