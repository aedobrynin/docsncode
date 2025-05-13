package paths

import "path/filepath"

func ConvertToPathInResultDir(pathToProjectRoot, target string, isFile bool, pathToResultDir string) (string, error) {
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
