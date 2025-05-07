package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"docsncode/buildcache"
	"docsncode/cfg"
	"docsncode/html"
	"docsncode/pathsignorer"
	"docsncode/utils"
)

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

func buildDocsncodeForFile(absPathToSourceFile, absPathToResultFile, absPathToResultDir, absPathToProjectRoot string) error {
	fileExtension := filepath.Ext(absPathToSourceFile)
	log.Printf("File extension is: %s", fileExtension)

	languageInfo, err := cfg.GetLanguageInfo(fileExtension)
	if err != nil {
		return fmt.Errorf("error on getting language info: %w", err)
	}
	log.Printf("Building html for %s", languageInfo.Language)

	file, err := os.Open(absPathToSourceFile)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %w", absPathToSourceFile, err)
	}
	defer file.Close()

	html, err := html.BuildHTML(file, *languageInfo, absPathToProjectRoot, absPathToSourceFile, absPathToResultDir, absPathToResultFile)
	if err != nil {
		return fmt.Errorf("error on bulding HTML for %s: %w", absPathToSourceFile, err)
	}

	resultFile, err := createFileAndNeededDirs(absPathToResultFile)
	if err != nil {
		return fmt.Errorf("couldn't create result file %s: %w", absPathToResultFile, err)
	}
	defer resultFile.Close()

	// TODO: писать сразу в файл с небольшим буффером?
	_, err = html.WriteTo(resultFile)
	if err != nil {
		return fmt.Errorf("error on writing HTML to stdout: %w", err)
	}
	return nil
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

	// TODO: параллельная обработка
	return filepath.WalkDir(pathToProjectRoot, func(path string, entry os.DirEntry, err error) error {
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

		relPathToEntry, err := filepath.Rel(pathToProjectRoot, absolutePathToEntry)
		if err != nil {
			log.Printf("error on building rel path to source file: %v", err)
			return nil
		}

		if entry.IsDir() {
			log.Printf("start walking through %s directory", path)

			// TODO: поправить RelPathFromProjectRoot
			if pathsIgnorer.ShouldIgnore(pathsignorer.RelPathFromProjectRoot(relPathToEntry)) {
				log.Printf("paths ignorer said to ignore the directory")
				return nil
			}
		}

		log.Printf("start building docsncode for %s", path)

		targetPath, err := utils.ConvertToPathInResultDir(pathToProjectRoot,
			path,
			true, // isFile
			pathToResultDir)
		if err != nil {
			log.Printf("error on building path to result file for %s: %v", path, err)
			return nil
		}

		// TODO: поправить RelPathFromProjectRoot
		if !buildCache.ShouldBuild(buildcache.RelPathFromProjectRoot(relPathToEntry)) {
			log.Printf("current result is actual according to build cache")
			return nil
		}

		// TODO: поправить RelPathFromProjectRoot
		if pathsIgnorer.ShouldIgnore(pathsignorer.RelPathFromProjectRoot(relPathToEntry)) {
			log.Printf("paths ignorer said to ignore the file")
			return nil
		}

		err = buildDocsncodeForFile(absolutePathToEntry, targetPath, pathToResultDir, pathToProjectRoot)
		if err != nil {
			log.Printf("error on building docsncode for %s: %v", path, err)
			return nil
		}
		// TODO: поправить RelPathFromProjectRoot
		buildCache.StoreBuildResult(buildcache.RelPathFromProjectRoot(relPathToEntry))

		return nil
	})
}
