package html

import (
	"log"
	"net/url"
	"path/filepath"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"docsncode/cfg"
	"docsncode/models"
	"docsncode/paths"
	"docsncode/pathsignorer"
)

type linksResolverTransformer struct {
	absPathToProjectRoot string
	absPathToCurrentFile string
	absPathToResultDir   string
	absPathToResultFile  string
	pathsIgnorer         pathsignorer.PathsIgnorer
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

func (t *linksResolverTransformer) willThereBeResultFileWithSuchPath(path models.RelPathFromProjectRoot) bool {
	if t.pathsIgnorer.ShouldIgnore(path) {
		log.Printf("path=%s is ignored by paths ignorer", path)
		return false
	}
	return cfg.GetLanguageNameIfSupported(filepath.Ext(string(path))) != nil
}

func (t *linksResolverTransformer) getUpdatedPath(path []byte) []byte {
	pathString := string(path)
	if isURL(pathString) {
		log.Printf("Destination is URL")
		return path
	}

	absPath := pathString
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(filepath.Dir(t.absPathToCurrentFile), absPath)
	} else {
		// TODO: make it a warning
		log.Println("found link with absolute path. It probably won't work on a different host")
	}

	log.Printf("absPath=%s", absPath)

	if !isPathNested(t.absPathToProjectRoot, absPath) {
		relPath, err := filepath.Rel(filepath.Dir(t.absPathToResultFile), absPath)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", t.absPathToResultFile, absPath, err)
			return path
		}
		return []byte(relPath)
	}

	log.Println("path is nested")

	relPathFromProjectRoot, err := filepath.Rel(t.absPathToProjectRoot, absPath)
	if err != nil {
		log.Printf("error on getting relative path for %s, %s: %s", t.absPathToProjectRoot, t.absPathToCurrentFile, err)
		return path
	}

	if t.willThereBeResultFileWithSuchPath(models.RelPathFromProjectRoot(relPathFromProjectRoot)) {
		log.Println("path will have result file")
		resultPath, err := paths.ConvertToPathInResultDir(t.absPathToProjectRoot, absPath, true, t.absPathToResultDir)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", t.absPathToResultFile, absPath, err)
			return path
		}

		relResultPath, err := filepath.Rel(filepath.Dir(t.absPathToResultFile), resultPath)
		if err != nil {
			log.Printf("error on getting relative path for %s, %s: %s", filepath.Dir(t.absPathToResultFile), resultPath, err)
			return path
		}
		return []byte(relResultPath)
	}

	relPath, err := filepath.Rel(t.absPathToResultDir, absPath)
	if err != nil {
		log.Printf("error on getting relative path for %s, %s: %s", t.absPathToResultFile, absPath, err)
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
			img.Destination = t.getUpdatedPath(img.Destination)
			log.Printf("Updated destination is %s", img.Destination)
			return ast.WalkContinue, nil
		}

		if node.Kind() == ast.KindLink {
			link := node.(*ast.Link)
			log.Printf("Found link with destination=%s", link.Destination)
			link.Destination = t.getUpdatedPath(link.Destination)
			log.Printf("Updated destination is %s", link.Destination)
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
}
