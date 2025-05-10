package buildcache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"docsncode/models"
)

type modificationTimeBasedBuildCache struct {
	absPathToProjectRoot   string
	absPathToCacheDataFile string
	absPathToResultDir     string

	previousCacheEntries map[models.RelPathFromProjectRoot]cacheEntry
	currentCacheEntries  sync.Map
}

type cacheEntry struct {
	SourceFileModTimestamp int64 `json:"source_file_modification_timestamp"`
	ResultFileModTimestamp int64 `json:"result_file_modification_timestamp"`
}

type cacheData struct {
	AbsPathToResultDir string                                       `json:"absolute_path_to_result_dir"`
	Entries            map[models.RelPathFromProjectRoot]cacheEntry `json:"entries"`
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

func getPreviousCacheEntries(absPathToCacheDataFile, absPathToResultDir string) map[models.RelPathFromProjectRoot]cacheEntry {
	file, err := os.Open(absPathToCacheDataFile)
	if os.IsNotExist(err) {
		log.Printf("There is no cache file with path %s", absPathToCacheDataFile)
		return make(map[models.RelPathFromProjectRoot]cacheEntry)
	}
	if err != nil {
		log.Printf("There is an error on opening cache data file with path %s: %s", absPathToCacheDataFile, err)
		return make(map[models.RelPathFromProjectRoot]cacheEntry)
	}
	defer file.Close()

	var previousCacheData cacheData
	err = json.NewDecoder(file).Decode(&previousCacheData)
	if err != nil {
		log.Printf("Error reading previous cache data path=%s, err= %s, will init empty cache", absPathToCacheDataFile, err)
		return make(map[models.RelPathFromProjectRoot]cacheEntry)
	} else {
		log.Printf("Successfully read previous cache data from file")
	}

	if previousCacheData.AbsPathToResultDir != absPathToResultDir {
		log.Printf("Cache data from file was built for different result dir, will init empty cache")
		return nil
	}

	return previousCacheData.Entries
}

func NewModificationTimeBasedBuildCache(absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile string) BuildCache {
	return &modificationTimeBasedBuildCache{
		absPathToProjectRoot:   absPathToProjectRoot,
		absPathToCacheDataFile: absPathToCacheDataFile,
		absPathToResultDir:     absPathToResultDir,
		previousCacheEntries:   getPreviousCacheEntries(absPathToCacheDataFile, absPathToResultDir),
		currentCacheEntries:    sync.Map{},
	}
}

// TODO: ок ли, что не возвращаем ошибки?
func (c *modificationTimeBasedBuildCache) ShouldBuild(relPathFromProjectRootToFile models.RelPathFromProjectRoot) bool {
	entry, isPresent := c.previousCacheEntries[relPathFromProjectRootToFile]
	if !isPresent {
		log.Printf("didn't find entry with path %s in cache", relPathFromProjectRootToFile)
		return true
	}

	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathFromProjectRootToFile))
	sourceFileModTimestamp := getModTimestamp(absPathToSourceFile)
	if sourceFileModTimestamp == nil {
		log.Printf("source file modification timestamp is nil")
		return true
	}
	if entry.SourceFileModTimestamp != *sourceFileModTimestamp {
		log.Printf("source file modification timestamp differs from the value saved in cache")
		return true
	}

	// TODO(important): заиспользовать utils.ConvertToPathInResultDir
	absPathToResultFile := filepath.Join(c.absPathToResultDir, string(relPathFromProjectRootToFile)+".html")
	resultFileModTimestamp := getModTimestamp(absPathToResultFile)
	if resultFileModTimestamp == nil {
		log.Printf("result file modification timestamp is nil")
		return true
	}
	if entry.ResultFileModTimestamp != *resultFileModTimestamp {
		log.Printf("result file modification timestamp differs from the value saved in cache")
		return true
	}

	c.currentCacheEntries.Store(relPathFromProjectRootToFile, entry)

	return false
}

// TODO: ок ли, что не возвращаем ошибки?
func (c *modificationTimeBasedBuildCache) StoreSuccessfulBuildResult(relPathFromProjectRootToFile models.RelPathFromProjectRoot) {
	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathFromProjectRootToFile))
	sourceFileModTimestamp := getModTimestamp(absPathToSourceFile)
	if sourceFileModTimestamp == nil {
		log.Printf("source file modification timestamp is nil, can't store it in cache")
		return
	}

	// TODO(important): заиспользовать utils.ConvertToPathInResultDir
	absPathToResultFile := filepath.Join(c.absPathToResultDir, string(relPathFromProjectRootToFile)+".html")
	resultFileModTimestamp := getModTimestamp(absPathToResultFile)
	if resultFileModTimestamp == nil {
		log.Printf("result file modification timestamp is nil, can't store it in cache")
		return
	}

	c.currentCacheEntries.Store(
		relPathFromProjectRootToFile,
		cacheEntry{
			SourceFileModTimestamp: *sourceFileModTimestamp,
			ResultFileModTimestamp: *resultFileModTimestamp,
		})
}

func (c *modificationTimeBasedBuildCache) Dump() error {
	entries := make(map[models.RelPathFromProjectRoot]cacheEntry)
	c.currentCacheEntries.Range(func(path any, entry any) bool {
		p, ok := path.(models.RelPathFromProjectRoot)
		if !ok {
			log.Fatalf("Unexpected key in current cache entries")
		}
		e, ok := entry.(cacheEntry)
		if !ok {
			log.Fatalf("Unexpected value in current cache entries")
		}
		entries[p] = e
		return true
	})
	cacheData := cacheData{
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
