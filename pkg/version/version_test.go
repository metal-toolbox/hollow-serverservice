package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	version = "1"
	commit = "abc123"
	date = "Some point in time"
	builtBy = "go test"

	assert.Equal(t, "1 (abc123@Some point in time by go test)", String())
}
