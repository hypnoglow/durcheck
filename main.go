package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"io"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

var (
	// definitions for testing purposes

	stderr io.Writer = os.Stderr
	exit             = os.Exit
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprint(os.Stderr, "\tdurcheck [package] ...\n")
		fmt.Fprint(os.Stderr, "\tdurcheck [directory] ...\n")
		fmt.Fprint(os.Stderr, "\tdurcheck ./... \n")
	}
	flag.Parse()

	importPaths := gotool.ImportPaths(flag.Args())
	if len(importPaths) == 0 {
		importPaths = []string{"."}
	}

	loadcfg := loader.Config{Build: &build.Default}

	unconsumed, err := loadcfg.FromArgs(importPaths, true)
	if err != nil {
		log.Fatalf("failed to parse arguments: %s", err)
	}
	if len(unconsumed) > 0 {
		log.Fatalf("unconsumed arguments: %v", unconsumed)
	}

	program, err := loadcfg.Load()
	if err != nil {
		log.Fatalf("failed to load packages: %s", err)
	}

	ntf := notifier{
		fset: program.Fset,
		out:  stderr,
	}

	var wasProblem bool
	for _, pkgInfo := range program.InitialPackages() {
		if len(pkgInfo.Files) == 0 {
			fmt.Fprintf(os.Stderr, "WARNING: no go files in package %q\n", pkgInfo.Pkg.Path())
			continue
		}

		insp := &inspector{tinf: pkgInfo.Info}
		for _, f := range pkgInfo.Files {
			problems := insp.inspect(f)
			for _, p := range problems {
				wasProblem = true
				ntf.notify(p)
			}
		}
	}

	code := 0
	if wasProblem {
		code = 1
	}
	exit(code)
}
