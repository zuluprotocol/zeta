package cmd_test

import (
	"testing"

	cmd "zuluprotocol/zeta/zeta/cmd/zetawallet/commands"
	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/flags"
	vgrand "zuluprotocol/zeta/zeta/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeAPITokenFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testDescribeAPITokenValidFlagsSucceeds)
	t.Run("Missing flags fails", testDescribeAPITokenWithMissingFlagsFails)
}

func testDescribeAPITokenValidFlagsSucceeds(t *testing.T) {
	// given
	testDir := t.TempDir()
	token := vgrand.RandomStr(10)
	_, passphraseFilePath := NewPassphraseFile(t, testDir)

	f := &cmd.DescribeAPITokenFlags{
		Token:          token,
		PassphraseFile: passphraseFilePath,
	}

	// when
	err := f.Validate()

	// then
	require.NoError(t, err)
}

func testDescribeAPITokenWithMissingFlagsFails(t *testing.T) {
	testDir := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, testDir)

	tcs := []struct {
		name        string
		flags       *cmd.DescribeAPITokenFlags
		missingFlag string
	}{
		{
			name: "without token",
			flags: &cmd.DescribeAPITokenFlags{
				PassphraseFile: passphraseFilePath,
			},
			missingFlag: "token",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// when
			err := tc.flags.Validate()

			// then
			assert.ErrorIs(t, err, flags.MustBeSpecifiedError(tc.missingFlag))
		})
	}
}
