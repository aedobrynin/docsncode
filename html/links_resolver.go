package html

import (
	"log"
	"net/url"
	"path/filepath"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"docsncode/cfg"
	"docsncode/utils"
)

type linksResolverTransformer struct {
	absPathToProjectRoot string
	absPathToCurrentFile string
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

func getUpdatedPath(path []byte, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string) []byte {
	pathString := string(path)
	if isURL(pathString) {
		log.Printf("Destination is URL")
		return path
	}

	absPath := pathString
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(filepath.Dir(absPathToCurrentFile), absPath)
	} else {
		// TODO: make it a warning
		log.Println("found link with absolute path. It probably won't work on a different host")
	}

	log.Printf("absPath=%s", absPath)

	if !isPathNested(absPathToProjectRoot, absPath) {
		relPath, err := filepath.Rel(filepath.Dir(absPathToResultFile), absPath)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", absPathToResultFile, absPath, err)
			return path
		}
		return []byte(relPath)
	}

	log.Println("path is nested")

	// TODO(important): переделать на нормальную функцию
	// TODO(important): учитывать pathsIgnorer
	_, err := cfg.GetLanguageInfo(filepath.Ext(absPath))
	if err == nil {
		log.Println("path will have result file")
		resultPath, err := utils.ConvertToPathInResultDir(absPathToProjectRoot, absPath, true, absPathToResultDir)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", absPathToResultFile, absPath, err)
			return path
		}

		relResultPath, err := filepath.Rel(filepath.Dir(absPathToResultFile), resultPath)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", filepath.Dir(absPathToResultFile), resultPath, err)
			return path
		}
		return []byte(relResultPath)
	}

	relPath, err := filepath.Rel(absPathToResultDir, absPath)
	if err != nil {
		log.Printf("error on getting relative path for %s, %s: %s", absPathToResultFile, absPath, err)
		return path
	}
	return []byte(relPath)
}

func (t *linksResolverTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if node.Kind() == ast.KindImage {
			img := node.(*ast.Image)
			log.Printf("Found image with destination=%s", img.Destination)
			img.Destination = getUpdatedPath(img.Destination, t.absPathToProjectRoot, t.absPathToCurrentFile, t.absPathToResultDir, t.absPathToResultFile)
			log.Printf("Updated destination is %s", img.Destination)
			return ast.WalkContinue, nil
		}

		if node.Kind() == ast.KindLink {
			link := node.(*ast.Link)
			log.Printf("Found link with destination=%s", link.Destination)
			link.Destination = getUpdatedPath(link.Destination, t.absPathToProjectRoot, t.absPathToCurrentFile, t.absPathToResultDir, t.absPathToResultFile)
			log.Printf("Updated destination is %s", link.Destination)
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
}
