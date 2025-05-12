package buildcache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"docsncode/models"
	"docsncode/paths"
)

type modificationTimeBasedCacheEntry struct {
	SourceFileModTimestamp int64 `json:"source_file_modification_timestamp"`
	ResultFileModTimestamp int64 `json:"result_file_modification_timestamp"`
}

type modificationTimeBasedCacheData = cacheData[modificationTimeBasedCacheEntry]

type modificationTimeBasedBuildCache struct {
	absPathToProjectRoot   string
	absPathToCacheDataFile string
	absPathToResultDir     string

	previousCacheEntries map[models.RelPathFromProjectRoot]modificationTimeBasedCacheEntry
	currentCacheEntries  sync.Map
}

func getModTimestamp(path string) *int64 {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("Didn't find file %s, will return nil mod timestamp", path)
		return nil
	} else if err != nil {
		log.Printf("Error on getting os.Stat for path=%s, err=%s, will return nil mod timestamp", path, err)
		return nil
	}

	modTimestamp := stat.ModTime().Unix()
	return &modTimestamp
}

func NewModificationTimeBasedBuildCache(absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile string) BuildCache {
	return &modificationTimeBasedBuildCache{
		absPathToProjectRoot:   absPathToProjectRoot,
		absPathToCacheDataFile: absPathToCacheDataFile,
		absPathToResultDir:     absPathToResultDir,
		previousCacheEntries:   getPreviousCacheEntries[modificationTimeBasedCacheEntry](absPathToCacheDataFile, absPathToResultDir),
		currentCacheEntries:    sync.Map{},
	}
}

func (c *modificationTimeBasedBuildCache) ShouldBuild(relPathToSourceFile models.RelPathFromProjectRoot) bool {
	entry, isPresent := c.previousCacheEntries[relPathToSourceFile]
	if !isPresent {
		log.Printf("didn't find entry with path %s in cache", relPathToSourceFile)
		return true
	}

	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathToSourceFile))
	sourceFileModTimestamp := getModTimestamp(absPathToSourceFile)
	if sourceFileModTimestamp == nil {
		log.Printf("source file modification timestamp is nil")
		return true
	}
	if entry.SourceFileModTimestamp != *sourceFileModTimestamp {
		log.Printf("source file modification timestamp differs from the value saved in cache")
		return true
	}

	// TODO(important): do not calculate path to result file here
	absPathToResultFile, err := paths.ConvertToPathInResultDir(c.absPathToProjectRoot, absPathToSourceFile, true, c.absPathToResultDir)
	if err != nil {
		log.Printf("Couldn't get abs path to result file: %s", err)
	}

	resultFileModTimestamp := getModTimestamp(absPathToResultFile)
	if resultFileModTimestamp == nil {
		log.Printf("result file modification timestamp is nil")
		return true
	}
	if entry.ResultFileModTimestamp != *resultFileModTimestamp {
		log.Printf("result file modification timestamp differs from the value saved in cache")
		return true
	}

	c.currentCacheEntries.Store(relPathToSourceFile, entry)

	return false
}

func (c *modificationTimeBasedBuildCache) StoreSuccessfulBuildResult(relPathToSourceFile models.RelPathFromProjectRoot, absPathToResultFile models.AbsPath) {
	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathToSourceFile))
	sourceFileModTimestamp := getModTimestamp(absPathToSourceFile)
	if sourceFileModTimestamp == nil {
		log.Printf("source file modification timestamp is nil, can't store it in cache")
		return
	}

	resultFileModTimestamp := getModTimestamp(string(absPathToResultFile))
	if resultFileModTimestamp == nil {
		log.Printf("result file modification timestamp is nil, can't store it in cache")
		return
	}

	c.currentCacheEntries.Store(
		relPathToSourceFile,
		modificationTimeBasedCacheEntry{
			SourceFileModTimestamp: *sourceFileModTimestamp,
			ResultFileModTimestamp: *resultFileModTimestamp,
		})
}

func (c *modificationTimeBasedBuildCache) Dump() error {
	entries := make(map[models.RelPathFromProjectRoot]modificationTimeBasedCacheEntry)
	c.currentCacheEntries.Range(func(path any, entry any) bool {
		p, ok := path.(models.RelPathFromProjectRoot)
		if !ok {
			log.Fatalf("Unexpected key in current cache entries")
		}
		e, ok := entry.(modificationTimeBasedCacheEntry)
		if !ok {
			log.Fatalf("Unexpected value in current cache entries")
		}
		entries[p] = e
		return true
	})
	cacheData := modificationTimeBasedCacheData{
		AbsPathToResultDir: c.absPathToResultDir,
		Entries:            entries,
	}

	file, err := os.Create(c.absPathToCacheDataFile)
	if err != nil {
		return fmt.Errorf("error on creating file for cache dump: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(cacheData)
	if err != nil {
		return fmt.Errorf("error on dumping cache to file: %w", err)
	}

	return nil
}
