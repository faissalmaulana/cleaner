package filepath

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilePaths(t *testing.T) {
	var roots = []string{".config", ".cache"}

	mockGetFilePath := func(path, pkg_name string) (string, error) {
		return filepath.Join(path, pkg_name), nil
	}

	t.Run("success get paths with valid input", func(t *testing.T) {
		var fakerootdir = "home/lizzy"
		input_pkg_name := "firefox"
		expected := []string{"home/lizzy/.config/firefox", "home/lizzy/.cache/firefox"}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		result, err := GetFilePaths(ctx, mockGetFilePath, fakerootdir, input_pkg_name, roots)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, result)
	})

	t.Run("success cancel search filepaths", func(t *testing.T) {
		var fakerootdir = "home/lizzy"
		input_pkg_name := "firefox"

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately — simulate user who already gave up

		result, err := GetFilePaths(ctx, mockGetFilePath, fakerootdir, input_pkg_name, roots)

		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, result)
	})

	t.Run("fail get paths because invalid input", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		result, err := GetFilePaths(ctx, mockGetFilePath, "", "", roots)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

}

func TestDeleteFilePaths(t *testing.T) {
	mockDeleteFilePathImplementation := func(path string) error {
		fmt.Printf("%s deleted\n", path)
		return nil
	}

	t.Run("success delete filepaths", func(t *testing.T) {

		filepaths := []string{".config/firefox", ".cache/firefox"}

		err := DeleteFilePaths(mockDeleteFilePathImplementation, filepaths)
		assert.NoError(t, err)
	})

	t.Run("partial fail delete filepaths", func(t *testing.T) {
		mockDeleteFilePathImplementation := func(path string) error {
			available := []string{".config/chrome"}

			found := slices.Contains(available, path)
			if found {
				fmt.Printf("deleted: %s", path)
				return nil
			} else {
				return errors.New("path not found")
			}

		}

		filepaths := []string{".config/chrome", ".cache/firefox"}

		err := DeleteFilePaths(mockDeleteFilePathImplementation, filepaths)
		assert.Error(t, err)
	})

	t.Run("fail delete filepaths because give empty input", func(t *testing.T) {
		err := DeleteFilePaths(mockDeleteFilePathImplementation, []string{})
		assert.Error(t, err)
	})

}

func TestGetFilePathFromOS(t *testing.T) {
	dirloc, err := setupTempDir("tmp", "hello_world.txt")

	if err != nil {
		t.Fatalf("something wrong: %v", err)
	}

	defer os.RemoveAll(dirloc)

	expected := filepath.Join(dirloc, "hello_world.txt")
	result, err := GetFilePathFromOS(dirloc, "hello", false)

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
