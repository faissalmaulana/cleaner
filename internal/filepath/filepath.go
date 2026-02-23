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

func GetFilePathFromOS(dirloc, pkg string) (string, error) {
	// read all files and folders inside the directory
	entries, err := os.ReadDir(dirloc)
	if err != nil {
		return "", err
	}

	// loop through each entry found in the directory
	for _, entry := range entries {
		if strings.Contains(entry.Name(), pkg) {
			fullpath := filepath.Join(dirloc, entry.Name())
			// combine the directory path with the exact file name
			fmt.Printf("find %s\n", fullpath)
			return fullpath, nil
		}
	}

	// if the loop finishes without finding anything, return an empty string
	return "", nil
}
