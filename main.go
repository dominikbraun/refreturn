package main

import (
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

func main() {
	for i := 0; i < numWorkers; i++ {
		gate.Add(1)
		go NewWorker().RecvTask(jobQueue, &gate)
	}

	if err := filepath.Walk(".", handler); err != nil {
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
