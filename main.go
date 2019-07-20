package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Printf("%s\n", "starting refreturn")

	if err := filepath.Walk(".", fileHandler); err != nil {
		log.Fatal(err)
	}
}

func fileHandler(path string, file os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	_ = readFile(path)

	return nil
}
