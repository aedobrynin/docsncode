package pathsignorer

import "docsncode/internal/models"

type AlwaysNotIgnoringPathsIgnorer struct {
}

func NewAlwaysNotIgnoringPathsIgnorer() PathsIgnorer {
	return &AlwaysNotIgnoringPathsIgnorer{}
}

func (*AlwaysNotIgnoringPathsIgnorer) ShouldIgnore(path models.RelPathFromProjectRoot) bool {
	return false
}
