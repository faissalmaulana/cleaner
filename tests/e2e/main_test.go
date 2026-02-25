package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Resolve project root (2 levels up from /tests/integration)
	rootDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		panic(err)
	}

	binDir := filepath.Join(rootDir, "bin")
	binaryPath = filepath.Join(binDir, "cleaner")

	if err := os.MkdirAll(binDir, 0755); err != nil {
		panic(err)
	}

	// Build coverage-instrumented binary
	cmd := exec.Command("go", "build", "-cover", "-o", binaryPath, "./main.go")
	cmd.Dir = rootDir

	if out, err := cmd.CombinedOutput(); err != nil {
		panic(string(out))
	}

	os.Exit(m.Run())
}

func setupFakeHome(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	return tmp
}

func runCmd(args ...string) ([]byte, error) {
	rootCoverDir := filepath.Join("..", "..", ".coverdata")

	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR="+rootCoverDir)

	return cmd.CombinedOutput()
}

func runCmdWithInput(input string, args ...string) ([]byte, error) {
	rootDir, _ := filepath.Abs(filepath.Join("..", ".."))
	coverDir := filepath.Join(rootDir, ".coverdata")

	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR="+coverDir)

	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	return cmd.CombinedOutput()
}

func createFakePkgName(pkgName string, xdgDirs ...string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	for _, xdgDir := range xdgDirs {
		fullPath := filepath.Join(homeDir, xdgDir, pkgName)

		err := os.MkdirAll(fullPath, 0o755)
		if err != nil {
			return err
		}
	}

	return nil
}
