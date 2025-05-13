package buildcache

import "docsncode/internal/models"

// TODO(important): учитывать версию утилиты в кэше?
// Например, если в новой версии утилиты меняется отображение для уже существующих элементов,
// 	то необходимо перегенерить результат

type BuildCache interface {
	// ShouldBuild and StoreBuildResult can be called concurrently
	// TODO: ок ли, что не возвращаем ошибки?
	ShouldBuild(relPathToSourceFile models.RelPathFromProjectRoot) bool
	// TODO: ок ли, что не возвращаем ошибки?
	StoreSuccessfulBuildResult(relPathToSourceFile models.RelPathFromProjectRoot, absPathToResultFile models.AbsPath)

	// Dump should be called not more than once.
	// The call must be after all ShouldBuild and StoreBuildResult calls.
	Dump() error
}
