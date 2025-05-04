package main

import (
	"fmt"
	"log"
	"os"
)

// @docsncode
// This is a comment block
//
// It can contain [links](https://example.com) and images ![cat](images/cat.png)
//
// This image is taken from this [file](images/cat.png)
//
// This [link](https://example.com "link with a title") has a title
// @docsncode

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

	// @docsncode
	// Here we use function from [html.go](html/html.go)
	// @docsncode
	err = buildDocsncode(pathToProjectRoot, pathToResultDir)
	if err != nil {
		log.Fatalf("error on building docsncode: %v", err)
	}
	log.Printf("written result to %s\n", pathToResultDir)

	/* @docsncode
	This is an example of multiline comment.

	It is

	very

	very

	very

	multiline.

	@docsncode */

}
