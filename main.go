package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: docsncode <path-to-file>")
		return
	}
	filePath := os.Args[1]

	fileExtension := filepath.Ext(filePath)
	if !IsExtensionSupported(fileExtension) {
		log.Fatalf("Files with extension %s are not supported", fileExtension)
	}
}
