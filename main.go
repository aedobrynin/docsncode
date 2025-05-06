package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"docsncode/buildcache"
)

// @docsncode
// This is a comment block
//
// It can contain [links](https://example.com) and images ![cat](images/cat.png)
//
// This image is taken from this [file](images/cat.png)
//
// This [link](https://example.com "link with a title") has a title
// @docsncode

func initBuildCache(forceRebuild, noCache bool, absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile string) buildcache.BuildCache {
	if noCache {
		log.Printf("will use always empty build cache")
		return buildcache.NewAlwaysEmptyBuildCache()
	}
	// TODO: handle force rebuild

	modificationTimeBasedBuildCache := buildcache.NewModificationTimeBasedBuildCache(absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile)
	if forceRebuild {
		log.Printf("will use force rebuild cache")
		return buildcache.NewForceRebuildCache(modificationTimeBasedBuildCache)
	}
	log.Printf("will use modification-time-based build cache")
	return modificationTimeBasedBuildCache
}

func main() {
	log.SetOutput(os.Stderr)

	// TOFIX: флаги не считываются
	forceRebuild := flag.Bool("force-rebuild", false, "ignore cached values and rebuild whole result")
	noCache := flag.Bool("no-cache", false, "ignore cached values and do not cache the result")

	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalf("Usage: docsncode <path-to-project-root> <path-to-result-dir> [path-to-cache-file] [--force-rebuild] [--no-cache]")
	}
	pathToProjectRoot := flag.Arg(0)
	pathToResultDir := flag.Arg(1)
	var pathToCacheFile string
	if flag.NArg() > 2 {
		pathToCacheFile = flag.Arg(2)
	} else {
		pathToCacheFile = filepath.Join(pathToProjectRoot, ".docsncode_cache.json")
	}

	log.Printf("path_to_project_root=%s, path_to_result_dir=%s, path_to_cache_file=%s, force_rebuild=%t, no_cache=%t", pathToProjectRoot, pathToResultDir, pathToCacheFile, *forceRebuild, *noCache)

	// TODO: правда ли, что это должно происходить тут?
	err := os.MkdirAll(pathToResultDir, 0755)
	if err != nil {
		log.Fatalf("error on creating result directory: %v", err)
	}

	absPathToProjectRoot, err := filepath.Abs(pathToProjectRoot)
	if err != nil {
		log.Fatalf("error on getting abs path to project root: %v", err)
	}

	absPathToResultDir, err := filepath.Abs(pathToResultDir)
	if err != nil {
		log.Fatalf("error on getting abs path to result dir: %v", err)
	}

	absPathToCacheDataFile, err := filepath.Abs(pathToCacheFile)
	if err != nil {
		log.Fatalf("error on getting abs path to cache data file: %v", err)
	}

	buildCache := initBuildCache(*forceRebuild, *noCache, absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile)

	// @docsncode
	// Here we use function from [html.go](html/html.go)
	// @docsncode
	err = buildDocsncode(pathToProjectRoot, pathToResultDir, buildCache)
	if err != nil {
		log.Fatalf("error on building docsncode: %v", err)
	}
	log.Printf("written result to %s", pathToResultDir)
	// TODO: не должны ли мы дампить кэш при ошибке?
	err = buildCache.Dump()
	if err != nil {
		log.Printf("error on dumping build cache: %v", err)
	}

	/* @docsncode
	This is an example of multiline comment.

	It is

	very

	very

	very

	multiline.

	@docsncode */

}
