package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

const usage = `Usage of %s:
	durcheck [flags] [package] ...
	durcheck [flags] [directory] ...
	durcheck [flags] ./..

Available flags:
	-t    Also check test packages
`

var (
	// definitions for testing purposes

	stderr io.Writer = os.Stderr
	exit             = os.Exit
)

func main() {
	lintTests := flag.Bool("t", false, "Also check test packages")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()

	importPaths := gotool.ImportPaths(flag.Args())
	if len(importPaths) == 0 {
		importPaths = []string{"."}
	}

	var wasProblem bool
	for _, p := range importPaths {
		loadcfg := loader.Config{Build: &build.Default}
		if *lintTests {
			loadcfg.ImportWithTests(p)
		} else {
			loadcfg.Import(p)
		}

		program, err := loadcfg.Load()
		if err != nil {
			log.Printf("failed to load packages: %s", err)
			continue
		}

		ntf := notifier{
			fset: program.Fset,
			out:  stderr,
		}

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
	}

	code := 0
	if wasProblem {
		code = 1
	}
	exit(code)
}
