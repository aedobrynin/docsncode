package buildcache

type alwaysEmptyBuildCache struct {
}

func NewAlwaysEmptyBuildCache() BuildCache {
	return &alwaysEmptyBuildCache{}
}

func (*alwaysEmptyBuildCache) ShouldBuild(relPathFromProjectRootToFile RelPathFromProjectRoot) bool {
	return true
}

func (*alwaysEmptyBuildCache) StoreBuildResult(relPathFromProjectRootToFile RelPathFromProjectRoot) {

}

func (*alwaysEmptyBuildCache) Dump() error {
	return nil
}
