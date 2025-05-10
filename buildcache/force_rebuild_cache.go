package buildcache

import "docsncode/models"

type ForceRebuildCache struct {
	// this cache will be used for StoreResult and Dump methods
	storingCache BuildCache
}

func NewForceRebuildCache(storingCache BuildCache) BuildCache {
	return &ForceRebuildCache{
		storingCache: storingCache,
	}
}

func (*ForceRebuildCache) ShouldBuild(relPathFromProjectRootToFile models.RelPathFromProjectRoot) bool {
	return true
}

func (c *ForceRebuildCache) StoreSuccessfulBuildResult(relPathFromProjectRootToFile models.RelPathFromProjectRoot) {
	c.storingCache.StoreSuccessfulBuildResult(relPathFromProjectRootToFile)
}

func (c *ForceRebuildCache) Dump() error {
	return c.storingCache.Dump()
}
