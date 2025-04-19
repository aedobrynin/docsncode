package main

import (
	"docsncode/cfg"
	"docsncode/html"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func buildDocsncodeForFile(path, resultPath string) error {
	// TODO: не запускать билд, если текущий результат актуален

	fileExtension := filepath.Ext(path)
	log.Printf("File extension is: %s\n", fileExtension)

	languageInfo, err := cfg.GetLanguageInfo(fileExtension)
	if err != nil {
		return fmt.Errorf("error on getting language info: %w", err)
	}
	log.Printf("Building html for %s\n", languageInfo.Language)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %w", path, err)
	}
	defer file.Close()

	html, err := html.BuildHTML(file, *languageInfo)
	if err != nil {
		return fmt.Errorf("error on bulding HTML for %s: %w", path, err)
	}

	resultFile, err := os.OpenFile(resultPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("couldn't create result file %s: %w", resultPath, err)
	}
	defer resultFile.Close()

	// TODO: писать сразу в файл с небольшим буффером?
	_, err = html.WriteTo(resultFile)
	if err != nil {
		return fmt.Errorf("error on writing HTML to stdout: %w", err)
	}
	return nil
}

func convertToPathInResultDir(pathToProjectRoot, target string, isFile bool, pathToResultDir string) (string, error) {
	relativePath, err := filepath.Rel(pathToProjectRoot, target)
	if err != nil {
		return "", err
	}

	result := filepath.Join(pathToResultDir, relativePath)
	if isFile {
		return result + ".html", nil
	}
	return result, nil
}

func buildDocsncode(pathToProjectRoot, pathToResultDir string) error {
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
			log.Printf("couldn't get absolute path of %s\n", path)
			return filepath.SkipDir
		}

		if pathToResultDir == absolutePathToEntry {
			log.Printf("skip walking through %s, because it matches result dir path", absolutePathToEntry)
			return filepath.SkipDir
		}

		if entry.IsDir() {
			log.Printf("start walking through %s directory\n", path)

			// TODO: не создавать директорию, если внутри неё нет файлов,
			// по которым построена документация

			targetPath, err := convertToPathInResultDir(
				pathToProjectRoot,
				path,
				false, // isFile
				pathToResultDir)
			if err != nil {
				log.Printf("error on building path to result dir for %s: %v", path, err)
				return nil
			}

			err = os.MkdirAll(targetPath, 0755)
			if err != nil {
				return fmt.Errorf("error on creating directory %s: %w", path, err)
			}
			return nil
		}

		log.Printf("start building docsncode for %s\n", path)

		targetPath, err := convertToPathInResultDir(pathToProjectRoot,
			path,
			true, // isFile
			pathToResultDir)
		if err != nil {
			log.Printf("error on building path to result file for %s: %v", path, err)
			return nil
		}

		err = buildDocsncodeForFile(path, targetPath)
		if err != nil {
			log.Printf("error on building docsncode for %s: %v", path, err)
			return nil
		}

		return nil
	})
}
