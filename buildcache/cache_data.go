package buildcache

import (
	"docsncode/models"
	"encoding/json"
	"log"
	"os"
)

type cacheData[cacheEntry any] struct {
	AbsPathToResultDir string                                       `json:"absolute_path_to_result_dir"`
	Entries            map[models.RelPathFromProjectRoot]cacheEntry `json:"entries"`
}

func getPreviousCacheEntries[cacheEntry any](absPathToCacheDataFile, absPathToResultDir string) map[models.RelPathFromProjectRoot]cacheEntry {
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

	var previousCacheData cacheData[cacheEntry]
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
