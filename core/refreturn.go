package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const numWorkers int = 10

var (
	jobQueue = make(chan string, 100)
	gate     = sync.WaitGroup{}
)

func Run(dir string) {
	for i := 0; i < numWorkers; i++ {
		gate.Add(1)
		go (&worker{}).RecvTask(jobQueue, &gate)
	}

	if err := filepath.Walk(dir, handler); err != nil {
		log.Fatal(err)
	}

	close(jobQueue)
	gate.Wait()
}

func handler(path string, file os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if strings.HasSuffix(path, ".go") {
		jobQueue <- path
	}
	return nil
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

	go func() {
		ast.Walk(v, file)
		close(v.idents)
	}()

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

		if v.checkRefs(fields) {
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
