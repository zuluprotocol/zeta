package paths_test

import (
	"testing"

	vgtest "zuluprotocol/zeta/zeta/libs/test"
	"zuluprotocol/zeta/zeta/paths"

	"github.com/stretchr/testify/assert"
)

func TestNewPaths(t *testing.T) {
	t.Run("Create a Paths without path returns the default implementation", testCreatingPathsWithoutPathReturnsDefaultImplementation)
	t.Run("Create a Paths without path returns the custom implementation", testCreatingPathsWithPathReturnsCustomImplementation)
}

func testCreatingPathsWithoutPathReturnsDefaultImplementation(t *testing.T) {
	p := paths.New("")

	assert.IsType(t, &paths.DefaultPaths{}, p)
}

func testCreatingPathsWithPathReturnsCustomImplementation(t *testing.T) {
	p := paths.New(vgtest.RandomPath())

	assert.IsType(t, &paths.CustomPaths{}, p)
}
