package pathsignorer

import (
	"fmt"
	"log"
	"os"

	ignore "github.com/sabhiram/go-gitignore"

	"docsncode/models"
)

type goGitignoreBasedPathsIgnorer struct {
	ignorer ignore.GitIgnore
}

func NewGoGitignoreBasedPathsIgnorer(pathToDocsncodeIgnore models.RelPathFromProjectRoot) (PathsIgnorer, error) {
	ignorer, err := ignore.CompileIgnoreFile(string(pathToDocsncodeIgnore))
	if os.IsNotExist(err) {
		log.Println("do not see .docsncodeignore file, will build empty paths ignorer")
		ignorer = ignore.CompileIgnoreLines()
	} else if err != nil {
		return nil, fmt.Errorf("error on builing go-gitignore ignorer: %w", err)
	}
	return &goGitignoreBasedPathsIgnorer{
		ignorer: *ignorer,
	}, nil
}

func (i *goGitignoreBasedPathsIgnorer) ShouldIgnore(path models.RelPathFromProjectRoot) bool {
	return i.ignorer.MatchesPath(string(path))
}
