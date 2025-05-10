package buildcache

import "docsncode/models"

type alwaysEmptyBuildCache struct {
}

func NewAlwaysEmptyBuildCache() BuildCache {
	return &alwaysEmptyBuildCache{}
}

func (*alwaysEmptyBuildCache) ShouldBuild(relPathFromProjectRootToFile models.RelPathFromProjectRoot) bool {
	return true
}

func (*alwaysEmptyBuildCache) StoreSuccessfulBuildResult(relPathFromProjectRootToFile models.RelPathFromProjectRoot) {

}

func (*alwaysEmptyBuildCache) Dump() error {
	return nil
}
