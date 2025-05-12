package pathsignorer

import "docsncode/models"

type PathsIgnorer interface {
	// Should be goroutine-safe
	ShouldIgnore(path models.RelPathFromProjectRoot) bool
}
