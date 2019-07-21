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
		fmt.Println(err)
		return nil
	}

	ast.Walk(visitor{}, f)

	return nil
}

type visitor struct {
	depth int
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		v.depth = 0
		return nil
	}

	tabs := strings.Repeat("\t", v.depth)
	fmt.Printf("%s%T\n", tabs, n)

	v.depth++
	return v
}
