package filepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// common xdg config directories
var dirlocs = []string{".config", ".cache"}

func GetFilePaths(fn func(path, pkg_name string) (string, error), homedir, pkg_name string) ([]string, error) {
	var errs []error
	findPaths := make([]string, 0)

	if err := validatePkgName(pkg_name); err != nil {
		return nil, err
	}

	fullpaths, err := combinePathsWithHomeDir(homedir, dirlocs)
	if err != nil {
		return nil, err
	}

	for _, dirloc := range fullpaths {
		path, err := fn(dirloc, pkg_name)
		if err != nil {
			errs = append(errs, err)
		}

		findPaths = append(findPaths, path)
	}

	return findPaths, errors.Join(errs...)
}

// validPkgRegex ensures the name contains only alphanumeric, spaces, dots, underscores, or dashes.
var validPkgRegex = regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)

func validatePkgName(pkg_name string) error {
	// check for empty or whitespace-only strings
	if strings.TrimSpace(pkg_name) == "" {
		return errors.New("pkg_name cannot be empty")
	}

	// use regex for character validation (blocks / \ and null bytes)
	if !validPkgRegex.MatchString(pkg_name) {
		return errors.New("pkg_name contains invalid characters")
	}

	// explicitly block path traversal '..' since regex allows single dots
	if strings.Contains(pkg_name, "..") {
		return errors.New("pkg_name cannot contain path traversal")
	}

	return nil
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
