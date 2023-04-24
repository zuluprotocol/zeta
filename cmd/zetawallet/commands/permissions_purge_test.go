package cmd_test

import (
	"testing"

	cmd "zuluprotocol/zeta/zeta/cmd/zetawallet/commands"
	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/flags"
	vgrand "zuluprotocol/zeta/zeta/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurgePermissionsFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testPurgePermissionsFlagsValidFlagsSucceeds)
	t.Run("Missing flags fails", testPurgePermissionsFlagsMissingFlagsFails)
}

func testPurgePermissionsFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	testDir := t.TempDir()
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	f := &cmd.PurgePermissionsFlags{
		Wallet:         vgrand.RandomStr(10),
		PassphraseFile: passphraseFilePath,
		Force:          true,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, f.Wallet, req.Wallet)
	assert.Equal(t, passphrase, req.Passphrase)
}

func testPurgePermissionsFlagsMissingFlagsFails(t *testing.T) {
	testDir := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, testDir)

	tcs := []struct {
		name        string
		flags       *cmd.PurgePermissionsFlags
		missingFlag string
	}{
		{
			name: "without wallet",
			flags: &cmd.PurgePermissionsFlags{
				Wallet:         "",
				PassphraseFile: passphraseFilePath,
			},
			missingFlag: "wallet",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// when
			req, err := tc.flags.Validate()

			// then
			assert.ErrorIs(t, err, flags.MustBeSpecifiedError(tc.missingFlag))
			require.Empty(t, req)
		})
	}
}
