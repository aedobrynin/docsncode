package main

import (
	"fmt"
	"log"
	"os"
)

// @docsncode_comment_block_start
// This is a comment block
//
// It can be multiline, contain [links](https://example.com) and images ![cat](../images/cat.png)
// @docsncode_comment_block_end

func main() {
	log.SetOutput(os.Stderr)

	// TODO: make path-to-result-dir optional
	if len(os.Args) < 3 {
		fmt.Println("Usage: docsncode <path-to-project-root> <path-to-result-dir>")
		return
	}
	pathToProjectRoot := os.Args[1]
	pathToResultDir := os.Args[2]

	err := os.MkdirAll(pathToResultDir, 0755)
	if err != nil {
		log.Fatalf("error on creating result directory: %v", err)
	}

	// @docsncode_comment_block_start
	// Here we use function from [html.go](html/html.go.html)
	// @docsncode_comment_block_end
	err = buildDocsncode(pathToProjectRoot, pathToResultDir)
	if err != nil {
		log.Fatalf("error on building docsncode: %v", err)
	}
	log.Printf("written result to %s\n", pathToResultDir)
}
