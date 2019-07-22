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

// Run spawns all workers and walking through the specified
// directory. Waits until all files have been processed.
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

// Handles any file that has been found by filepath.Walk().
// If the file is a Go source file, it will be put into the queue.
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

// RecvTask pops tasks off the global job queue and processes
// each job. Stops when no more jobs are in the queue.
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

// process parses a source file and walks through its AST in a
// seperate goroutine, receiving the affected idents from it.
func (w *worker) process(path string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil
	}

	idents := make(chan match)
	v := vtor{idents}

	go func() {
		ast.Walk(v, file)
		close(v.idents)
	}()

	for {
		match, ok := <-idents
		if !ok {
			break
		}
		pos := fset.PositionFor(match.pos, false)
		fn := match.idt.Name

		fmt.Printf("%s: %s\n", pos, fn)
	}
	return nil
}

// vtor satisfies the ast.Visitor interface and is used by for
// inspecting every AST node in ast.Walk().
type vtor struct {
	idents chan<- match
}

// match represents a match, i. e. a reference-returning func.
// pos describes the function's position in the source file, while
// idt holds the actual name.
type match struct {
	pos token.Pos
	idt *ast.Ident
}

// Visit checks the type a given AST node `n`. If the node is a
// function declaration, a return type check is performed.
func (v vtor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch decl := n.(type) {
	case *ast.FuncDecl:
		fields := decl.Type.Results

		if v.checkRefs(fields) {
			v.idents <- match{decl.Pos(), decl.Name}
		}
	}
	return v
}

// checkRefs determines if one of a function's return types
// is a reference. Each return type is an entry in `fields`.
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
