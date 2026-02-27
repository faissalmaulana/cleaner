package e2e

import (
	"testing"

	"github.com/faissalmaulana/cleaner/cmd"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {

	actual, err := runCmd("--version")
	expected := cmd.Version + "\n"

	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}
