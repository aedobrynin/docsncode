package buildcache

import "docsncode/models"

// TODO(important): учитывать версию утилиты в кэше?
// Например, если в новой версии утилиты меняется отображение для уже существующих элементов,
// 	то необходимо перегенерить результат

type BuildCache interface {
	// ShouldBuild and StoreBuildResult can be called concurrently
	ShouldBuild(relPathFromProjectRootToFile models.RelPathFromProjectRoot) bool
	StoreSuccessfulBuildResult(relPathFromProjectRootToFile models.RelPathFromProjectRoot)

	// Dump should be called not more than once.
	// The call must be after all ShouldBuild and StoreBuildResult calls.
	Dump() error
}
