package cmd_test

import (
	"testing"

	cmd "zuluprotocol/zeta/zeta/cmd/zetawallet/commands"
	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/flags"
	vgrand "zuluprotocol/zeta/zeta/libs/rand"
	"zuluprotocol/zeta/zeta/wallet/service/v2/connections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAPITokenFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testGenerateAPITokenValidFlagsSucceeds)
	t.Run("Missing flags fails", testGenerateAPITokenWithMissingFlagsFails)
}

func testGenerateAPITokenValidFlagsSucceeds(t *testing.T) {
	// given
	testDir := t.TempDir()
	description := vgrand.RandomStr(10)
	wallet := vgrand.RandomStr(10)
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletPassphrase, walletPassphraseFilePath := NewPassphraseFile(t, testDir)

	f := &cmd.GenerateAPITokenFlags{
		Description:          description,
		PassphraseFile:       passphraseFilePath,
		WalletName:           wallet,
		WalletPassphraseFile: walletPassphraseFilePath,
	}

	expectedReq := connections.GenerateAPITokenParams{
		Description: description,
		Wallet: connections.GenerateAPITokenWalletParams{
			Name:       wallet,
			Passphrase: walletPassphrase,
		},
	}
	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.EqualValues(t, expectedReq, req)
}

func testGenerateAPITokenWithMissingFlagsFails(t *testing.T) {
	testDir := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	_, walletPassphraseFilePath := NewPassphraseFile(t, testDir)

	tcs := []struct {
		name        string
		flags       *cmd.GenerateAPITokenFlags
		missingFlag string
	}{
		{
			name: "without wallet name",
			flags: &cmd.GenerateAPITokenFlags{
				PassphraseFile:       passphraseFilePath,
				WalletPassphraseFile: walletPassphraseFilePath,
			},
			missingFlag: "wallet-name",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// when
			req, err := tc.flags.Validate()

			// then
			assert.ErrorIs(t, err, flags.MustBeSpecifiedError(tc.missingFlag))
			assert.Empty(t, req)
		})
	}
}
