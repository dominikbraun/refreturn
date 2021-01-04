package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	numWorkers    int    = 4
	fileExtension string = ".go"
)

// Run spawns all workers and walking through the specified
// directory. Waits until all files have been processed.
func Run(dir string) error {
	var wg sync.WaitGroup
	jobQueue := make(chan string)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			worker := &Worker{}
			_ = worker.readFromQueue(jobQueue, &wg)
		}()
	}

	if err := sendFiles(dir, jobQueue); err != nil {
		return err
	}

	close(jobQueue)
	wg.Wait()

	return nil
}

func sendFiles(dir string, jobQueue chan<- string) error {
	return filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, fileExtension) {
			jobQueue <- path
		}
		return nil
	})
}

type Worker struct{}

// readFromQueue pops tasks off the global job queue and processes
// each job. Stops when no more jobs are in the queue.
func (w *Worker) readFromQueue(jobs <-chan string, wg *sync.WaitGroup) error {
	for path := range jobs {
		if err := w.findAllocationsInFile(path); err != nil {
			return err
		}
	}
	wg.Done()

	return nil
}

// findAllocationsInFile parses a source file and walks through its AST in a
// seperate goroutine, receiving the affected matches from it.
func (w *Worker) findAllocationsInFile(path string) error {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		return nil
	}

	visitor := Visitor{
		matches: make(chan Node),
	}

	go func() {
		ast.Walk(visitor, file)
		close(visitor.matches)
	}()

	for match := range visitor.matches {
		pos := fileSet.PositionFor(match.Position, false)
		fn := match.Identifier.Name

		fmt.Printf("%s: %s\n", pos, fn)
	}

	return nil
}

// Visitor satisfies the ast.Visitor interface and is used by for
// inspecting every AST node in ast.Walk().
type Visitor struct {
	matches chan Node
	filter  func(node ast.Node) bool
}

// Node represents a Node, i. e. a reference-returning func.
// Position describes the function's Position in the source file, while
// Identifier holds the actual name.
type Node struct {
	Position   token.Pos
	Identifier *ast.Ident
}

// Visit checks the type a given AST node `n`. If the node is a
// function declaration, a return type check is performed.
func (v Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch decl := node.(type) {
	case *ast.FuncDecl:
		if containsReference(decl.Type.Results) {
			v.matches <- Node{
				Position:   decl.Pos(),
				Identifier: decl.Name,
			}
		}
	}
	return v
}

// containsReference determines if one of a function's return types
// is a reference. Each return type is an entry in `fields`.
func containsReference(fieldList *ast.FieldList) bool {
	if fieldList == nil {
		return false
	}
	for _, f := range fieldList.List {
		if _, ok := f.Type.(*ast.StarExpr); ok {
			return true
		}
	}
	return false
}
