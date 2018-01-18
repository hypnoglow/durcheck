package main

import (
	"fmt"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
)

type notifier struct {
	fset *token.FileSet
	out  io.Writer
}

func (n notifier) notify(p problem) {
	pos := n.fset.Position(p.call.Pos())
	expr := types.ExprString(p.call)
	msg := fmt.Sprintf("implicit time.Duration means nanoseconds in %q ", expr)
	fmt.Fprintf(n.out, "%s:%d:%d: %s\n", shortPath(pos.Filename), pos.Line, pos.Column, msg)
}

func shortPath(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}

	if rel, err := filepath.Rel(cwd, path); err == nil && len(rel) < len(path) {
		return rel
	}

	return path
}
