package buildcache

import "docsncode/models"

type alwaysEmptyBuildCache struct {
}

func NewAlwaysEmptyBuildCache() BuildCache {
	return &alwaysEmptyBuildCache{}
}

func (*alwaysEmptyBuildCache) ShouldBuild(relPathToSourceFile models.RelPathFromProjectRoot) bool {
	return true
}

func (*alwaysEmptyBuildCache) StoreSuccessfulBuildResult(relPathToSourceFile models.RelPathFromProjectRoot, absPathToResultFile models.AbsPath) {

}

func (*alwaysEmptyBuildCache) Dump() error {
	return nil
}
