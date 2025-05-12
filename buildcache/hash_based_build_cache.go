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
	"docsncode/paths"
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

func (c *hashBasedBuildCache) ShouldBuild(relPathToSourceFile models.RelPathFromProjectRoot) bool {
	entry, isPresent := c.previousCacheEntries[relPathToSourceFile]
	if !isPresent {
		log.Printf("didn't find entry with path %s in cache", relPathToSourceFile)
		return true
	}

	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathToSourceFile))
	sourceFileHash, err := calculateSHA256(absPathToSourceFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of source file, err=%s", err)
		return true
	}
	if entry.SourceFileHash != sourceFileHash {
		log.Printf("source file hash differs from the value saved in cache")
		return true
	}

	// TODO(important): do not calculate path to result file here
	absPathToResultFile, err := paths.ConvertToPathInResultDir(c.absPathToProjectRoot, absPathToSourceFile, true, c.absPathToResultDir)
	if err != nil {
		log.Printf("Couldn't get abs path to result file: %s", err)
	}

	resultFileHash, err := calculateSHA256(absPathToResultFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of result file, err=%s", err)
		return true
	}
	if entry.ResultFileHash != resultFileHash {
		log.Printf("result file hash differs from the value saved in cache")
		return true
	}

	c.currentCacheEntries.Store(relPathToSourceFile, entry)

	return false
}

func (c *hashBasedBuildCache) StoreSuccessfulBuildResult(relPathToSourceFile models.RelPathFromProjectRoot, absPathToResultFile models.AbsPath) {
	absPathToSourceFile := filepath.Join(c.absPathToProjectRoot, string(relPathToSourceFile))
	sourceFileHash, err := calculateSHA256(absPathToSourceFile)
	if err != nil {
		log.Printf("Couldn't calculate hash of source file, err=%s. Can't store it in cache", err)
		return
	}

	resultFileHash, err := calculateSHA256(string(absPathToResultFile))
	if err != nil {
		log.Printf("Couldn't calculate hash of result file, err=%s. Can't store it in cache", err)
		return
	}

	c.currentCacheEntries.Store(
		relPathToSourceFile,
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
