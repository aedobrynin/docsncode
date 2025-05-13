package pathsignorer

import "docsncode/internal/models"

type PathsIgnorer interface {
	// Should be goroutine-safe
	ShouldIgnore(path models.RelPathFromProjectRoot) bool
}
