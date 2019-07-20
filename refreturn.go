package main

import "os"

func readFile(path string) error {
	_, err := os.Open(path)
	return err
}

func returnTypes(signature string) []string {
	return []string{}
}
