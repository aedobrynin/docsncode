package pathsignorer

// TODO: move to models
type RelPathFromProjectRoot string

type PathsIgnorer interface {
	ShouldIgnore(path RelPathFromProjectRoot) bool
}
