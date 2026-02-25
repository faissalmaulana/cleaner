package utils

import (
	"errors"
	"regexp"
	"strings"
)

// validPkgRegex ensures the name contains only alphanumeric, spaces, dots, underscores, or dashes.
var validPkgRegex = regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)

func ValidatePkgName(pkg_name string) error {
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
