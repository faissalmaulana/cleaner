package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	out, err := runCmd("--version")

	assert.NoError(t, err)
	assert.Equal(t, "0.0.0\n", string(out))

}
