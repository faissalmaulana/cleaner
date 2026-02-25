package integrations_test

import (
	"os"
	"path/filepath"
	"testing"

	fp "github.com/faissalmaulana/cleaner/internal/filepath"
	"github.com/stretchr/testify/assert"
)

func TestGetFilePaths(t *testing.T) {
	fakehomedir, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(fakehomedir)
	})

	config := filepath.Join(fakehomedir, ".config")
	configpkgname := filepath.Join(config, "firefox")

	cache := filepath.Join(fakehomedir, ".cache")
	cachepkgname := filepath.Join(cache, "firefox")

	err = os.MkdirAll(configpkgname, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(cachepkgname, 0755)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("success get filepaths", func(t *testing.T) {
		getFromOS := func(path, pkgName string) (string, error) {
			return fp.GetFilePathFromOS(path, pkgName, false)
		}
		input := "firefox"
		expected := []string{configpkgname, cachepkgname}

		result, err := fp.GetFilePaths(getFromOS, fakehomedir, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("fail get filepaths because pkg_name not found", func(t *testing.T) {
		getFromOS := func(path, pkgName string) (string, error) {
			return fp.GetFilePathFromOS(path, pkgName, false)
		}

		input := "somethingnotthere"
		result, err := fp.GetFilePaths(getFromOS, fakehomedir, input)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeleteFilePaths(t *testing.T) {
	fakehomedir, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(fakehomedir)
	})

	config := filepath.Join(fakehomedir, ".config")
	configpkgname := filepath.Join(config, "firefox")

	cache := filepath.Join(fakehomedir, ".cache")
	cachepkgname := filepath.Join(cache, "firefox")

	err = os.MkdirAll(configpkgname, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(cachepkgname, 0755)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("success delete filepaths", func(t *testing.T) {
		samplefilepaths := []string{configpkgname, cachepkgname}
		err = fp.DeleteFilePaths(os.Remove, samplefilepaths)
		assert.NoError(t, err)
	})

	t.Run("fail delete filepaths because input filepaths is not found", func(t *testing.T) {
		samplefilepaths := []string{filepath.Join(config, "chrome"), filepath.Join(cache, "chrome")}
		err = fp.DeleteFilePaths(os.Remove, samplefilepaths)
		assert.Error(t, err)
	})
}
