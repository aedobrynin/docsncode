package buildcache

type ForceRebuildCache struct {
	// this cache will be used for StoreResult and Dump methods
	storingCache BuildCache
}

func NewForceRebuildCache(storingCache BuildCache) BuildCache {
	return ForceRebuildCache{
		storingCache: storingCache,
	}
}

func (ForceRebuildCache) ShouldBuild(relPathFromProjectRootToFile RelPathFromProjectRoot) bool {
	return true
}

func (c ForceRebuildCache) StoreBuildResult(relPathFromProjectRootToFile RelPathFromProjectRoot) {
	c.storingCache.StoreBuildResult(relPathFromProjectRootToFile)
}

func (c ForceRebuildCache) Dump() error {
	return c.storingCache.Dump()
}
