package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUninstallSubCommand(t *testing.T) {
	tests := []struct {
		name        string
		pkgName     string
		query       string
		xdgDirs     []string
		withExact   bool
		userInput   string
		expectExist bool
	}{
		{
			name:        "fuzzy match removes package",
			pkgName:     "my-app",
			query:       "my",
			xdgDirs:     []string{".cache", ".config"},
			withExact:   false,
			userInput:   "y\n",
			expectExist: false,
		},
		{
			name:        "exact match removes package",
			pkgName:     "go",
			query:       "go",
			xdgDirs:     []string{".cache", ".config"},
			withExact:   true,
			userInput:   "y\n",
			expectExist: false,
		},
		{
			name:        "user aborts uninstall",
			pkgName:     "abort-app",
			query:       "abort",
			xdgDirs:     []string{".cache", ".config"},
			withExact:   false,
			userInput:   "n\n",
			expectExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := setupFakeHome(t)

			err := createFakePkgName(tt.pkgName, tt.xdgDirs...)
			require.NoError(t, err)

			args := []string{"uninstall"}
			if tt.withExact {
				args = append(args, "--ex")
			}
			args = append(args, tt.query)

			out, err := runCmdWithInput(tt.userInput, args...)
			require.NoError(t, err, "output: %s", string(out))

			for _, dir := range tt.xdgDirs {
				target := filepath.Join(home, dir, tt.pkgName)

				_, statErr := os.Stat(target)
				exists := statErr == nil

				if tt.expectExist {
					assert.True(t, exists, "expected %s to exist", target)
				} else {
					assert.False(t, exists, "expected %s to be removed", target)
				}
			}

			// Assert prompt behavior
			if tt.userInput == "n\n" {
				assert.Contains(t, string(out), "Aborted.")
			}
		})
	}
}

func TestUninstallSubCommandPackageNotFound(t *testing.T) {
	setupFakeHome(t)
	// create fake an existing package
	err := createFakePkgName("go", ".config", ".cache")
	if err != nil {
		t.Fatal(err)
	}

	// input a non-existing package
	input := "php"

	_, err = runCmdWithInput("y\n", "uninstall", input)
	assert.Error(t, err)

}
