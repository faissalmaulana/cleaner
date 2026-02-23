package filepath

import (
	"errors"
	"fmt"
)

// common xdg config directories
var dirlocs = []string{".config", ".cache"}

func GetFilePaths(fn func(path, pkg_name string) string, pkg_name string) []string {
	findPaths := make([]string, 0)

	// TODO: resolve the dirlocs path before passed to fn
	for _, dirloc := range dirlocs {
		findPaths = append(findPaths, fn(dirloc, pkg_name))
	}

	return findPaths
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
