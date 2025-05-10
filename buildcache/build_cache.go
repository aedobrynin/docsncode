package buildcache

// TODO: move to models and use in other places?
type RelPathFromProjectRoot string

// TODO: учитывать версию утилиты в кэше?
// Например, если в новой версии утилиты меняется отображение для уже существующих элементов,
// 	то необходимо перегенерить результат

type BuildCache interface {
	// ShouldBuild and StoreBuildResult can be called concurrently
	ShouldBuild(relPathFromProjectRootToFile RelPathFromProjectRoot) bool
	// TODO: rename if only one argument
	StoreBuildResult(relPathFromProjectRootToFile RelPathFromProjectRoot)

	// Dump should be called not more than once.
	// The call must be after all ShouldBuild and StoreBuildResult calls.
	Dump() error
}
