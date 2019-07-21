package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"sync"
)

type Worker interface {
	RecvTask(<-chan string, *sync.WaitGroup)
	process(string) error
}

func NewWorker() Worker {
	w := worker{}
	return &w
}

type worker struct{}

func (w *worker) RecvTask(jobs <-chan string, gate *sync.WaitGroup) {
	for {
		path, ok := <-jobs
		if !ok {
			break
		}
		w.process(path)
	}
	gate.Done()
}

func (w *worker) process(path string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil
	}

	idents := make(chan *ast.Ident)
	v := vtor{idents}

	go ast.Walk(v, file)

	for {
		i, ok := <-idents
		if !ok {
			break
		}
		fmt.Printf("%s returns one or more references.\n", i.Name)
	}
	return nil
}

type vtor struct {
	idents chan<- *ast.Ident
}

func (v vtor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch decl := n.(type) {
	case *ast.FuncDecl:
		fields := decl.Type.Results
		hasRefs := v.checkRefs(fields)

		if hasRefs {
			v.idents <- decl.Name
		}
	}
	return v
}

func (vtor) checkRefs(fields *ast.FieldList) bool {
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
