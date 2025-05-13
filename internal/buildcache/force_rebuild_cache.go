package buildcache

import "docsncode/internal/models"

type ForceRebuildCache struct {
	// this cache will be used for StoreResult and Dump methods
	storingCache BuildCache
}

func NewForceRebuildCache(storingCache BuildCache) BuildCache {
	return &ForceRebuildCache{
		storingCache: storingCache,
	}
}

func (*ForceRebuildCache) ShouldBuild(relPathToSourceFile models.RelPathFromProjectRoot) bool {
	return true
}

func (c *ForceRebuildCache) StoreSuccessfulBuildResult(relPathToSourceFile models.RelPathFromProjectRoot, absPathToResultFile models.AbsPath) {
	c.storingCache.StoreSuccessfulBuildResult(relPathToSourceFile, absPathToResultFile)
}

func (c *ForceRebuildCache) Dump() error {
	return c.storingCache.Dump()
}
