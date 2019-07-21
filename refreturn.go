package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func processFile(path string) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil
	}

	ast.Walk(vtor{}, f)

	return nil
}

type vtor struct{}

func (v vtor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.FuncDecl:
		fields := d.Type.Results

		if v.hasRefs(fields) {
			fmt.Printf("%v returns one or more references.\n", d.Name)
		}
	}
	return v
}

func (vtor) hasRefs(fields *ast.FieldList) bool {
	if fields == nil {
		return false
	}
	for _, f := range fields.List {
		if _, ok := f.Type.(*ast.StarExpr); ok {
			return true
		}
	}
	return false
}

func f1() (r *strings.Reader) {
	r = strings.NewReader("")
	return
}

func f2() *strings.Reader {
	return strings.NewReader("")
}

func f3() (int, int) {
	return 0, 0
}
