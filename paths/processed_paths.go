package paths

import (
	"path/filepath"
	"sync"

	"docsncode/models"
)

// Parallel calls to Update is goroutine-safe
// All IsFileProcessed calls and IsDirProcessed calls must be after all Update calls
type ProcessedPaths struct {
	processedFiles map[models.RelPathFromResultDir]struct{}
	processedDirs  map[models.RelPathFromResultDir]struct{}
	mut            sync.Mutex
}

func NewProcessedPaths() *ProcessedPaths {
	return &ProcessedPaths{
		processedFiles: make(map[models.RelPathFromResultDir]struct{}),
		processedDirs:  make(map[models.RelPathFromResultDir]struct{}),
		mut:            sync.Mutex{},
	}
}

func (pp *ProcessedPaths) IsFileProcessed(relPathToFile models.RelPathFromResultDir) bool {
	_, exists := pp.processedFiles[relPathToFile]
	return exists
}

func (pp *ProcessedPaths) IsDirProcessed(relPathToDir models.RelPathFromResultDir) bool {
	_, exists := pp.processedDirs[relPathToDir]
	return exists
}

func (pp *ProcessedPaths) Update(relPathToFile models.RelPathFromResultDir) {
	pp.mut.Lock()
	defer pp.mut.Unlock()

	pp.processedFiles[relPathToFile] = struct{}{}
	relPath := filepath.Dir(string(relPathToFile))
	for relPath != "." {
		pp.processedDirs[models.RelPathFromResultDir(relPath)] = struct{}{}
		relPath = filepath.Dir(relPath)
	}
}
