package pathsignorer

type AlwaysNotIgnoringPathsIgnorer struct {
}

func NewAlwaysNotIgnoringPathsIgnorer() PathsIgnorer {
	return &AlwaysNotIgnoringPathsIgnorer{}
}

func (*AlwaysNotIgnoringPathsIgnorer) ShouldIgnore(path RelPathFromProjectRoot) bool {
	return false
}
