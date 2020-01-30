package main

import (
	"go/ast"
	"go/types"
)

type inspector struct {
	tinf     types.Info
	problems []problem
}

func (insp *inspector) addProblem(p problem) {
	insp.problems = append(insp.problems, p)
}

func (insp *inspector) cleanProblems() {
	insp.problems = nil
}

func (insp *inspector) inspect(n ast.Node) []problem {
	defer insp.cleanProblems()
	ast.Inspect(n, insp.inspectFunc)
	return insp.problems
}

func (insp *inspector) inspectFunc(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.CallExpr:
		sel, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		fn, ok := insp.tinf.ObjectOf(sel.Sel).(*types.Func)
		if !ok {
			break
		}

		sig, ok := fn.Type().(*types.Signature)
		if !ok {
			break
		}

		for i := 0; i < sig.Params().Len(); i++ {
			param := sig.Params().At(i)
			if param.Type().String() != "time.Duration" {
				continue
			}

			if insp.isSuspicious(n.Args[i]) {
				p := problem{
					call: n,
				}
				insp.addProblem(p)
			}
		}
	}
	return true
}

func (insp *inspector) isSuspicious(param ast.Expr) bool {
	switch p := param.(type) {

	case *ast.BasicLit:
		//if p.Kind != token.INT {
		//	return false
		//}
		if p.Value == "0" {
			// 0 is unambiguous
			return false
		}
		return true

	case *ast.BinaryExpr:
		typ := insp.tinf.TypeOf(p.X)
		if typ.String() == "time.Duration" {
			return false
		}

		return true

	case *ast.Ident:
		obj := insp.tinf.ObjectOf(p)
		c, ok := obj.(*types.Const)
		if !ok {
			return false
		}

		b, ok := c.Type().(*types.Basic)
		if !ok {
			return false
		}

		if b.Kind() == types.UntypedInt {
			return true
		}
	}

	return false // assume everything else is ok
}
