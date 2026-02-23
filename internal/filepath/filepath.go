package filepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// common xdg config directories
var dirlocs = []string{".config", ".cache"}

func GetFilePaths(fn func(path, pkg_name string) string, homedir, pkg_name string) ([]string, error) {
	findPaths := make([]string, 0)

	fullpaths, err := combinePathsWithHomeDir(homedir, dirlocs)
	if err != nil {
		return nil, err
	}

	for _, dirloc := range fullpaths {
		findPaths = append(findPaths, fn(dirloc, pkg_name))
	}

	return findPaths, nil
}

func combinePathsWithHomeDir(homedir string, paths []string) ([]string, error) {
	var fullapaths []string

	for _, path := range paths {
		fullapaths = append(fullapaths, filepath.Join(homedir, path))
	}

	return fullapaths, nil
}

func DeleteFilePaths(fn func(path string) error, filepaths []string) error {
	var errs []error

	for _, filepath := range filepaths {
		err := fn(filepath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete %s: %w", filepath, err))
		}
	}

	return errors.Join(errs...)
}

func getFilePathFromOS(dirloc, pkg string, isExact bool) (string, error) {
	entries, err := os.ReadDir(dirloc)
	if err != nil {
		return "", err
	}

	// pre-lower the search term once
	searchLower := strings.ToLower(pkg)

	for _, entry := range entries {
		name := entry.Name()

		if isExact && name == pkg {
			return returnPath(dirloc, name)
		}

		if !isExact && strings.Contains(strings.ToLower(name), searchLower) {
			return returnPath(dirloc, name)
		}
	}

	return "", fmt.Errorf("not found")
}

func returnPath(dir, name string) (string, error) {
	fullPath := filepath.Join(dir, name)
	fmt.Printf("find %s\n", fullPath)
	return fullPath, nil
}
