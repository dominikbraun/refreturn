package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Printf("%s\n", "starting refreturn")

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
