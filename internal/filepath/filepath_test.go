package filepath

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilePaths(t *testing.T) {
	tests := []struct {
		name      string
		homedir   string
		pkgName   string
		implFn    func(path, pkgName string) (string, error)
		expected  []string
		expectErr bool
		errMsg    string
	}{
		{
			name:    "basic package with .config and .cache",
			homedir: "home/foo",
			pkgName: "firefox",
			implFn: func(path, pkgName string) (string, error) {
				return filepath.Join(path, pkgName), nil
			},
			expected: []string{
				filepath.Join("home/foo", ".config/firefox"),
				filepath.Join("home/foo", ".cache/firefox"),
			},
			expectErr: false,
		},
		{
			name:    "package name with spaces",
			homedir: "/home/user",
			pkgName: "visual studio code",
			implFn: func(path, pkgName string) (string, error) {
				return filepath.Join(path, pkgName), nil
			},
			expected: []string{
				filepath.Join("/home/user", ".config/visual studio code"),
				filepath.Join("/home/user", ".cache/visual studio code"),
			},
			expectErr: false,
		},
		{
			name:      "error: empty package name",
			homedir:   "/home/user",
			pkgName:   "",
			implFn:    nil,
			expectErr: true,
			errMsg:    "pkg_name cannot be empty",
		},
		{
			name:      "error: path traversal",
			homedir:   "/home/user",
			pkgName:   "../../etc/passwd",
			implFn:    nil,
			expectErr: true,
			errMsg:    "invalid characters", // Regex catches slashes first
		},
		{
			name:      "error: directory dots",
			homedir:   "/home/user",
			pkgName:   "invalid..name",
			implFn:    nil,
			expectErr: true,
			errMsg:    "pkg_name cannot contain path traversal",
		},
		{
			name:      "error: invalid characters (null byte)",
			homedir:   "/home/user",
			pkgName:   "bad\x00pkg",
			implFn:    nil,
			expectErr: true,
			errMsg:    "invalid characters",
		},
		{
			name:      "error: invalid characters (backslash)",
			homedir:   "/home/user",
			pkgName:   "wrong\\path",
			implFn:    nil,
			expectErr: true,
			errMsg:    "invalid characters",
		},
		{
			name:    "absolute path home directory",
			homedir: "/root",
			pkgName: "slack",
			implFn: func(path, pkgName string) (string, error) {
				return filepath.Join(path, pkgName), nil
			},
			expected: []string{
				filepath.Join("/root", ".config/slack"),
				filepath.Join("/root", ".cache/slack"),
			},
			expectErr: false,
		},
		{
			name:    "implFn returns custom path format",
			homedir: "/home/test",
			pkgName: "node",
			implFn: func(path, pkgName string) (string, error) {
				return fmt.Sprintf("%s/%s/config", path, pkgName), nil
			},
			expected: []string{
				"/home/test/.config/node/config",
				"/home/test/.cache/node/config",
			},
			expectErr: false,
		},
		{
			name:    "implFn returns error",
			homedir: "/home/user",
			pkgName: "test",
			implFn: func(path, pkgName string) (string, error) {
				return "", fmt.Errorf("path not found")
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:    "single character package name",
			homedir: "/home/user",
			pkgName: "a",
			implFn: func(path, pkgName string) (string, error) {
				return filepath.Join(path, pkgName), nil
			},
			expected: []string{
				filepath.Join("/home/user", ".config/a"),
				filepath.Join("/home/user", ".cache/a"),
			},
			expectErr: false,
		},
		{
			name:    "package name with special characters",
			homedir: "/home/user",
			pkgName: "my-app_v1.0",
			implFn: func(path, pkgName string) (string, error) {
				return filepath.Join(path, pkgName), nil
			},
			expected: []string{
				filepath.Join("/home/user", ".config/my-app_v1.0"),
				filepath.Join("/home/user", ".cache/my-app_v1.0"),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFilePaths(tt.implFn, tt.homedir, tt.pkgName)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeleteFilePaths(t *testing.T) {
	mockDeleteFilePathImplementation := func(path string) error {
		fmt.Printf("%s deleted\n", path)
		return nil
	}

	filepaths := []string{".config/firefox", ".cache/firefox"}
	err := DeleteFilePaths(mockDeleteFilePathImplementation, filepaths)

	assert.NoError(t, err, "should success delete filepaths")
}

func TestGetFilePathFromOS(t *testing.T) {
	dirloc, err := setupTempDir("tmp", "hello_world.txt")

	if err != nil {
		t.Fatalf("something wrong: %v", err)
	}

	defer os.RemoveAll(dirloc)

	expected := filepath.Join(dirloc, "hello_world.txt")
	result, err := getFilePathFromOS(dirloc, "hello", false)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

}

func TestCombinePathsWithHomeDir(t *testing.T) {
	var fakehomedir = "home/foo"
	paths := []string{"foo", "bar", filepath.Join("something", "woo"), "hello/world"}
	expected := []string{filepath.Join(fakehomedir, "foo"), filepath.Join(fakehomedir, "bar"), filepath.Join(fakehomedir, "something/woo"), filepath.Join(fakehomedir, "hello/world")}

	result, err := combinePathsWithHomeDir(fakehomedir, paths)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func setupTempDir(dirname string, filename string) (string, error) {
	dname, err := os.MkdirTemp("", dirname)
	if err != nil {
		return "", err
	}

	fname := filepath.Join(dname, filename)

	err = os.WriteFile(fname, []byte("hello world"), 0666)
	if err != nil {
		// return the directory path anyway so the caller can clean it up if needed.
		return dname, err
	}

	return dname, nil
}
