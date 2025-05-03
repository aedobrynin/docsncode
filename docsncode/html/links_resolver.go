package html

import (
	"docsncode/cfg"
	"docsncode/utils"
	"log"
	"net/url"
	"path/filepath"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type linksResolverTransformer struct {
	absPathToProjectRoot string
	absPathToResultDir   string
	absPathToResultFile  string
}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}

func isPathNested(parentPath, childPath string) bool {
	parentAbsPath, err := filepath.Abs(parentPath)
	if err != nil {
		log.Printf("error on getting absolute parent path: %v", err)
		return false
	}

	childAbsPath, err := filepath.Abs(childPath)
	if err != nil {
		log.Printf("error on getting absolute child path: %v", err)
		return false
	}

	return filepath.IsAbs(childAbsPath) && filepath.HasPrefix(childAbsPath, parentAbsPath)
}

func getUpdatedPath(path []byte, absPathToProjectRoot, absPathToResultDir, absPathToResultFile string) []byte {
	pathString := string(path)
	if isURL(pathString) {
		log.Printf("Destination is URL")
		return path
	}

	absPath, err := filepath.Abs(pathString)
	if err != nil {
		log.Printf("error on getting absolute path: %s; won't update the path", err)
		return path
	}

	log.Printf("absPath=%s\n", absPath)

	if !isPathNested(absPathToProjectRoot, absPath) {
		relPath, err := filepath.Rel(absPathToResultFile, absPath)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", absPathToResultFile, absPath, err)
			return path
		}
		return []byte(relPath)
	}

	log.Println("path is nested")

	// TODO: переделать на нормальную функцию
	_, err = cfg.GetLanguageInfo(filepath.Ext(absPath))
	if err == nil {
		log.Println("path will have result file")
		resultPath, err := utils.ConvertToPathInResultDir(absPathToProjectRoot, absPath, true, absPathToResultDir)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", absPathToResultFile, absPath, err)
			return path
		}
		return []byte(resultPath)
	}

	relPath, err := filepath.Rel(absPathToResultDir, absPath)
	if err != nil {
		log.Printf("error on getting relative path for %s, %s: %s", absPathToResultFile, absPath, err)
		return path
	}
	return []byte(relPath)
}

func (t *linksResolverTransformer) traverseChildren(node ast.Node) {
	if node == nil {
		return
	}

	if node.HasChildren() {
		t.traverseChildren(node.FirstChild())
	}
	t.traverseChildren(node.NextSibling())

	if node.Kind() == ast.KindImage {
		img := node.(*ast.Image)
		log.Printf("Found image with destination=%s", img.Destination)
		img.Destination = getUpdatedPath(img.Destination, t.absPathToProjectRoot, t.absPathToResultDir, t.absPathToResultFile)
		log.Printf("Updated destination is %s", img.Destination)
		return
	}

	if node.Kind() == ast.KindLink {
		link := node.(*ast.Link)
		log.Printf("Found link with destination=%s", link.Destination)
		link.Destination = getUpdatedPath(link.Destination, t.absPathToProjectRoot, t.absPathToResultDir, t.absPathToResultFile)
		log.Printf("Updated destination is %s", link.Destination)
		return
	}
}

func (t *linksResolverTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	t.traverseChildren(node)
}
