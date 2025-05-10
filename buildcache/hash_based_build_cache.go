package buildcache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"docsncode/models"
)

type hashBasedCacheEntry struct {
	SourceFileHash string `json:"source_file_hash"`
	ResultFileHash string `json:"result_file_hash"`
}

type hashBasedCacheData = cacheData[hashBasedCacheEntry]

type hashBasedBuildCache struct {
	absPathToProjectRoot   string
	absPathToCacheDataFile string
	absPathToResultDir     string

	previousCacheEntries map[models.RelPathFromProjectRoot]hashBasedCacheEntry
	currentCacheEntries  sync.Map
}

func calculateSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func NewHashBasedBuildCache(absPathToProjectRoot, absPathToResultDir, absPathToCacheDataFile string) BuildCache {
	return &hashBasedBuildCache{
		absPathToProjectRoot:   absPathToProjectRoot,
		absPathToCacheDataFile: absPathToCacheDataFile,
		absPathToResultDir:     absPathToResultDir,
		previousCacheEntries:   getPreviousCacheEntries[hashBasedCacheEntry](absPathToCacheDataFile, absPathToResultDir),
		currentCacheEntries:    sync.Map{},
	}
}

func (c *hashBasedBuildCache) ShouldBuild(relPathFromProjectRootToFile models.RelPathFromProjectRoot) bool {
	entry, isPresent := c.previousCacheEntries[relPathFromProjectRootToFile]
	if !isPresent {
		log.Printf("didn't find entry with path %s in cache", relPathFromProjectRootToFile)
		return true
	}

	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathFromProjectRootToFile))
	sourceFileHash, err := calculateSHA256(absPathToSourceFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of source file, err=%s", err)
		return true
	}
	if entry.SourceFileHash != sourceFileHash {
		log.Printf("source file hash differs from the value saved in cache")
		return true
	}

	// TODO(important): заиспользовать utils.ConvertToPathInResultDir
	absPathToResultFile := filepath.Join(c.absPathToResultDir, string(relPathFromProjectRootToFile)+".html")
	resultFileHash, err := calculateSHA256(absPathToResultFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of result file, err=%s", err)
		return true
	}
	if entry.ResultFileHash != resultFileHash {
		log.Printf("result file hash differs from the value saved in cache")
		return true
	}

	c.currentCacheEntries.Store(relPathFromProjectRootToFile, entry)

	return false
}

func (c *hashBasedBuildCache) StoreSuccessfulBuildResult(relPathFromProjectRootToFile models.RelPathFromProjectRoot) {
	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathFromProjectRootToFile))
	sourceFileHash, err := calculateSHA256(absPathToSourceFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of source file, err=%s. Can't store it in cache", err)
		return
	}

	// TODO(important): заиспользовать utils.ConvertToPathInResultDir
	absPathToResultFile := filepath.Join(c.absPathToResultDir, string(relPathFromProjectRootToFile)+".html")
	resultFileHash, err := calculateSHA256(absPathToResultFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of result file, err=%s. Can't store it in cache", err)
		return
	}

	c.currentCacheEntries.Store(
		relPathFromProjectRootToFile,
		hashBasedCacheEntry{
			SourceFileHash: sourceFileHash,
			ResultFileHash: resultFileHash,
		})
}

func (c *hashBasedBuildCache) Dump() error {
	entries := make(map[models.RelPathFromProjectRoot]hashBasedCacheEntry)
	c.currentCacheEntries.Range(func(path any, entry any) bool {
		p, ok := path.(models.RelPathFromProjectRoot)
		if !ok {
			log.Fatalf("Unexpected key in current cache entries")
		}
		e, ok := entry.(hashBasedCacheEntry)
		if !ok {
			log.Fatalf("Unexpected value in current cache entries")
		}
		entries[p] = e
		return true
	})
	cacheData := hashBasedCacheData{
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
