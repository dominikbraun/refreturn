package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := filepath.Walk(".", handler); err != nil {
		log.Fatal(err)
	}
}

func handler(path string, file os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if strings.HasSuffix(path, ".go") {
		_ = processFile(path)
	}

	return nil
}
