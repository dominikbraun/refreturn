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

// Run executes refreturn. It is in charge of spawning all worker
// routines, sending the files to be processed through a channel
// and gracefully stopping all workers.
//
// At the moment, the results are not queued but printed directly
// by the workers instead. This may change in the future.
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

// sendFiles sends all files matching the configured file extension
// through the jobQueue channel.
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

// readFromQueue pops off items from the job queue sequentially.
// Each queue item will be passed to findAllocationsInFile.
func (w *Worker) readFromQueue(jobs <-chan string, wg *sync.WaitGroup) error {
	for path := range jobs {
		if err := w.findAllocationsInFile(path); err != nil {
			return err
		}
	}
	wg.Done()

	return nil
}

// findAllocationsInFile parses a source code file and walks through
// its syntax tree. In doing so, a custom node visitor will check if
// a node is a function that returns a pointer. All matching nodes
// will be sent to a dedicated channel.
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

	// Iterate over all matches and print the file information for
	// each match. This should be outsourced to an own component.
	for match := range visitor.matches {
		pos := fileSet.PositionFor(match.Position, false)
		fn := match.Identifier.Name

		fmt.Printf("%s: %s\n", pos, fn)
	}

	return nil
}

// Visitor satisfies the ast.Visitor interface and is used by for
// inspecting every AST node using ast.Walk().
type Visitor struct {
	matches chan Node
	filter  func(node ast.Node) bool
}

// Node represents an AST node with an identifier.
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
// is a reference. Each return type is an entry in the FieldList.
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
