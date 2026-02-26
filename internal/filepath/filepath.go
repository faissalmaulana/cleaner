package filepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type result struct {
	path string
	err  error
}

func GetFilePaths(
	fn func(path, pkgName string) (string, error),
	homedir, pkgName string,
	roots []string,
) ([]string, error) {
	if homedir == "" || pkgName == "" {
		return nil, errors.New("homedir and pkgName is required")
	}

	fullpaths, err := combinePathsWithHomeDir(homedir, roots)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make(chan result)

	for _, dirloc := range fullpaths {
		dir := dirloc
		wg.Add(1)

		go func() {
			defer wg.Done()
			path, err := fn(dir, pkgName)
			results <- result{path, err}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var (
		findPaths []string
		errs      []error
	)

	for r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		}

		if r.path != "" {
			findPaths = append(findPaths, r.path)
		}
	}

	return findPaths, errors.Join(errs...)
}

func combinePathsWithHomeDir(homedir string, paths []string) ([]string, error) {
	var fullapaths []string

	for _, path := range paths {
		fullapaths = append(fullapaths, filepath.Join(homedir, path))
	}

	return fullapaths, nil
}

func DeleteFilePaths(fn func(path string) error, filepaths []string) error {
	if len(filepaths) == 0 || filepaths == nil {
		return errors.New("filepaths is required")
	}

	var errs []error

	for _, filepath := range filepaths {
		err := fn(filepath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete %s: %w", filepath, err))
		}
	}

	return errors.Join(errs...)
}

func GetFilePathFromOS(dirloc, pkg string, isExact bool) (string, error) {
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
