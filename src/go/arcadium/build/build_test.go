package build_test

import (
	"testing"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/build"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
)

func TestString(t *testing.T) {
	info := build.Info("Testing", "Version", "Branch", "Commit", "Date")
	info.Go = "Go"
	assert.Equal(t, info.String(), "Testing Version (branch: Branch, commit: Commit, date: Date, go: Go)")
}
