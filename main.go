package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const numWorkers int = 10

var (
	jobQueue = make(chan string, 100)
	done     = make(chan bool)
)

func main() {
	for i := 0; i < numWorkers; i++ {
		go NewWorker().RecvTask(jobQueue, done)
	}

	if err := filepath.Walk(".", handler); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numWorkers; i++ {
		<-done
	}
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
