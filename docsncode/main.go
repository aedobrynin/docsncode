package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"docsncode/cfg"
	"docsncode/html"
)

// @docsncode_comment_block_start
// This is a comment block
// It can be multiline
// @docsncode_comment_block_end

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) < 2 {
		fmt.Println("Usage: docsncode <path-to-file>")
		return
	}
	filePath := os.Args[1]

	fileExtension := filepath.Ext(filePath)

	log.Printf("File extension is: %s\n", fileExtension)

	languageInfo, err := cfg.GetLanguageInfo(fileExtension)
	if err != nil {
		log.Fatalf("error on getting language info: %v", err)
	}
	log.Printf("Building html for %s\n", languageInfo.Language)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("couldn't open file: %v", err)
	}
	defer file.Close()

	html, err := html.BuildHTML(file, *languageInfo)
	if err != nil {
		log.Fatalf("error on bulding HTML: %v", err)
	}

	_, err = html.WriteTo(os.Stdout)
	if err != nil {
		log.Fatalf("error on writing HTML to stdout: %v", err)
	}
}
