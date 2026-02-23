package filepath

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilePaths(t *testing.T) {
	mockGetFilePathImplementation := func(pathname, pkg_name string) string {
		return fmt.Sprintf("%s/%s", pathname, pkg_name)
	}

	pkg_name := "firefox"
	expected := []string{".config/firefox", ".cache/firefox"}

	result := GetFilePaths(mockGetFilePathImplementation, pkg_name)
	assert.Equal(t, expected, result, "should be equal")

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
	result, err := GetFilePathFromOS(dirloc, "hello")

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
