package main

import (
	"log"
)

var version = "UNDEFINED"

func main() {
	if err := rootCommand(version).Execute(); err != nil {
		log.Fatal(err)
	}
}
