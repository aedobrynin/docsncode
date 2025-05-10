package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"docsncode/buildcache"
	"docsncode/models"
	"docsncode/pathsignorer"
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

	cmd := &cli.Command{
		Name:  "docsncode",
		Usage: "An application to unite code and documentation",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force-rebuild",
				Usage: "Ignore cached result and build new result",
			},
			&cli.BoolFlag{
				Name:  "no-cache",
				Usage: "Ignore cached result and do not cache new result",
			},
		},
		UsageText: "docsncode <path-to-project-root> <path-to-result-dir> [path-to-cache-file] [--force-rebuild] [--no-cache]",
		Action: func(_ context.Context, c *cli.Command) error {
			if c.Args().Len() < 1 {
				log.Fatal("path-to-project-root is not provided")
			}
			if c.Args().Len() < 2 {
				log.Fatal("path-to-result-dir is not provided")
			}
			if c.Args().Len() > 3 {
				log.Fatal("Too many positional args")
			}
			pathToProjectRoot := c.Args().Get(0)
			pathToResultDir := c.Args().Get(1)
			pathToCacheFile := c.Args().Get(2)
			if pathToCacheFile == "" {
				pathToCacheFile = filepath.Join(pathToProjectRoot, ".docsncode_cache.json")
			}
			forceRebuild := c.Bool("force-rebuild")
			noCache := c.Bool("no-cache")

			log.Printf("path_to_project_root=%s, path_to_result_dir=%s, path_to_cache_file=%s, force_rebuild=%t, no_cache=%t", pathToProjectRoot, pathToResultDir, pathToCacheFile, forceRebuild, noCache)

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

			buildCache := initBuildCache(forceRebuild, noCache, absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile)

			// @docsncode
			// Here we use function from [html.go](html/html.go)
			// @docsncode

			pathToDocsncodeIgnoreFile := models.RelPathFromProjectRoot(filepath.Join(absPathToProjectRoot, ".docsncodeignore"))
			pathsIgnorer, err := pathsignorer.NewGoGitignoreBasedPathsIgnorer(pathToDocsncodeIgnoreFile)
			if err != nil {
				log.Fatalf("error on building paths ignorer: %v", err)
			}

			err = buildDocsncode(pathToProjectRoot, pathToResultDir, buildCache, pathsIgnorer)
			if err != nil {
				log.Fatalf("error on building docsncode: %v", err)
			}
			log.Printf("written result to %s", pathToResultDir)
			// TODO: не должны ли мы дампить кэш при ошибке?
			err = buildCache.Dump()
			if err != nil {
				log.Printf("error on dumping build cache: %v", err)
			}

			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
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
