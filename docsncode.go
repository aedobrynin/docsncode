package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"docsncode/buildcache"
	"docsncode/cfg"
	"docsncode/html"
	"docsncode/models"
	"docsncode/paths"
	"docsncode/pathsignorer"
)

var ErrLanguageNotSupported = errors.New("language is not supported")

func createFileAndNeededDirs(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func buildDocsncodeForFile(absPathToProjectRoot, absPathToSourceFile, absPathToResultDir, absPathToResultFile string, pathsIgnorer pathsignorer.PathsIgnorer) error {
	fileExtension := filepath.Ext(absPathToSourceFile)
	log.Printf("File extension is: %s", fileExtension)

	language := cfg.GetLanguageNameIfSupported(fileExtension)
	if language == nil {
		return ErrLanguageNotSupported
	}
	log.Printf("Building html for %s", *language)

	file, err := os.Open(absPathToSourceFile)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %w", absPathToSourceFile, err)
	}
	defer file.Close()

	html, err := html.BuildHTML(file, *language, absPathToProjectRoot, absPathToSourceFile, absPathToResultDir, absPathToResultFile, pathsIgnorer)
	if err != nil {
		return fmt.Errorf("error on bulding HTML for %s: %w", absPathToSourceFile, err)
	}

	resultFile, err := createFileAndNeededDirs(absPathToResultFile)
	if err != nil {
		return fmt.Errorf("couldn't create result file %s: %w", absPathToResultFile, err)
	}
	defer resultFile.Close()

	// TODO: писать сразу в файл с небольшим буффером?
	_, err = resultFile.Write(html)
	if err != nil {
		return fmt.Errorf("error on writing HTML to file: %w", err)
	}
	return nil
}

type buildTask struct {
	absPathToProjectRoot string
	absPathToSourceFile  string
	absPathToResultDir   string
	absPathToResultFile  string
	relPathToSourceFile  models.RelPathFromProjectRoot
}

func pushBuildTasks(tasksChan chan<- buildTask, pathToProjectRoot, pathToResultDir string, buildCache buildcache.BuildCache, pathsIgnorer pathsignorer.PathsIgnorer) {
	defer close(tasksChan)
	filepath.WalkDir(pathToProjectRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			log.Printf("error on opening %s: %v", path, err)
			return err
		}

		absolutePathToEntry, err := filepath.Abs(path)
		if err != nil {
			log.Printf("couldn't get absolute path of %s", path)
			return filepath.SkipDir
		}

		if pathToResultDir == absolutePathToEntry {
			log.Printf("skip walking through %s, because it matches result dir path", absolutePathToEntry)
			return filepath.SkipDir
		}

		var relPathToEntry models.RelPathFromProjectRoot
		{
			relPath, err := filepath.Rel(pathToProjectRoot, absolutePathToEntry)
			if err != nil {
				log.Printf("error on building rel path to source file: %v", err)
				return nil
			}
			relPathToEntry = models.RelPathFromProjectRoot(relPath)
		}

		if entry.IsDir() {
			log.Printf("start walking through %s directory", path)
			if pathsIgnorer.ShouldIgnore(relPathToEntry) {
				log.Printf("paths ignorer said to ignore the directory")
				return filepath.SkipDir
			}
			return nil
		}

		log.Printf("start building docsncode for %s", path)

		targetPath, err := paths.ConvertToPathInResultDir(pathToProjectRoot,
			path,
			true, // isFile
			pathToResultDir)
		if err != nil {
			log.Printf("error on building path to result file for %s: %v", path, err)
			return nil
		}

		if !buildCache.ShouldBuild(relPathToEntry) {
			log.Printf("current result is actual according to build cache")
			return nil
		}

		if pathsIgnorer.ShouldIgnore(relPathToEntry) {
			log.Printf("paths ignorer said to ignore the file")
			return nil
		}

		tasksChan <- buildTask{
			absPathToProjectRoot: pathToProjectRoot,
			absPathToSourceFile:  absolutePathToEntry,
			absPathToResultDir:   pathToResultDir,
			absPathToResultFile:  targetPath,
			relPathToSourceFile:  relPathToEntry,
		}

		log.Printf("pushed build task for path %s", absolutePathToEntry)
		return nil
	})
}

func processTasks(tasksChan <-chan buildTask, buildCache buildcache.BuildCache, pathsIgnorer pathsignorer.PathsIgnorer) *paths.ProcessedPaths {
	wg := sync.WaitGroup{}
	processedPaths := paths.NewProcessedPaths()

	for task := range tasksChan {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := buildDocsncodeForFile(task.absPathToProjectRoot, task.absPathToSourceFile, task.absPathToResultDir, task.absPathToResultFile, pathsIgnorer)
			if err != nil {
				log.Printf("Error on building result for path=%s, err=%s", task.relPathToSourceFile, err)
			} else {
				relPathToResultFile, err := filepath.Rel(task.absPathToResultDir, task.absPathToResultFile)
				if err != nil {
					log.Printf("error on getting relative path from %s to %s: %s", task.absPathToResultDir, task.absPathToResultFile, err)
					return
				}
				processedPaths.Update(models.RelPathFromResultDir(relPathToResultFile))
				buildCache.StoreSuccessfulBuildResult(task.relPathToSourceFile, models.AbsPath(task.absPathToResultFile))
			}
		}()
	}

	wg.Wait()
	return processedPaths
}

func removeUnrelatedPaths(pathToResultDir string, processedPaths *paths.ProcessedPaths) {
	filepath.WalkDir(pathToResultDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			log.Printf("error on opening %s: %v", path, err)
			return err
		}

		absolutePathToEntry, err := filepath.Abs(path)
		if err != nil {
			log.Printf("couldn't get absolute path of %s", path)
			return filepath.SkipDir
		}

		if pathToResultDir == absolutePathToEntry {
			return nil
		}

		var relPathToEntry models.RelPathFromResultDir
		{
			relPath, err := filepath.Rel(pathToResultDir, absolutePathToEntry)
			if err != nil {
				log.Printf("error on building rel path to source file: %v", err)
				return nil
			}
			relPathToEntry = models.RelPathFromResultDir(relPath)
		}

		if (entry.IsDir() && !processedPaths.IsDirProcessed(relPathToEntry)) ||
			(!entry.IsDir() && !processedPaths.IsFileProcessed(relPathToEntry)) {
			os.RemoveAll(absolutePathToEntry)
			log.Printf("Deleted file %s, because it's not supposed to be in the result directory", relPathToEntry)
		}
		return nil
	})
}

func buildDocsncode(pathToProjectRoot, pathToResultDir string, buildCache buildcache.BuildCache, pathsIgnorer pathsignorer.PathsIgnorer) error {
	pathToProjectRoot, err := filepath.Abs(pathToProjectRoot)
	if err != nil {
		return fmt.Errorf("couldn't get absolute path for project root directory: %w", err)
	}

	pathToResultDir, err = filepath.Abs(pathToResultDir)
	if err != nil {
		return fmt.Errorf("couldn't get absolute path for result directory: %w", err)
	}

	buildTasks := make(chan buildTask, 1)

	go pushBuildTasks(buildTasks, pathToProjectRoot, pathToResultDir, buildCache, pathsIgnorer)
	processedPaths := processTasks(buildTasks, buildCache, pathsIgnorer)
	removeUnrelatedPaths(pathToResultDir, processedPaths)
	return nil
}
