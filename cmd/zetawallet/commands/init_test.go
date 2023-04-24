package cmd_test

import (
	"testing"

	cmd "code.zetaprotocol.io/vega/cmd/vegawallet/commands"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("Initialising software succeeds", testInitialisingSoftwareSucceeds)
	t.Run("Forcing software initialisation succeeds", testForcingSoftwareInitialisationSucceeds)
}

func testInitialisingSoftwareSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := &cmd.InitFlags{
		Force: false,
	}

	// when
	err := cmd.Init(testDir, f)

	// then
	require.NoError(t, err)
}

func testForcingSoftwareInitialisationSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := &cmd.InitFlags{
		Force: false,
	}

	// when
	err := cmd.Init(testDir, f)

	// then
	require.NoError(t, err)

	// given
	f = &cmd.InitFlags{
		Force: true,
	}

	// when
	err = cmd.Init(testDir, f)

	// then
	require.NoError(t, err)
}
