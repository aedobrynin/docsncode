package pathsignorer

import "docsncode/models"

type PathsIgnorer interface {
	ShouldIgnore(path models.RelPathFromProjectRoot) bool
}
