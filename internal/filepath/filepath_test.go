package filepath

import (
	"fmt"
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
