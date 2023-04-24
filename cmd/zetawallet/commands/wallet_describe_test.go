package cmd_test

import (
	"testing"

	cmd "zuluprotocol/zeta/zeta/cmd/zetawallet/commands"
	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/flags"
	vgrand "zuluprotocol/zeta/zeta/libs/rand"
	"zuluprotocol/zeta/zeta/wallet/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWalletInfoFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testGetWalletInfoFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testGetWalletInfoFlagsMissingWalletFails)
}

func testGetWalletInfoFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.DescribeWalletFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := api.AdminDescribeWalletParams{
		Wallet:     walletName,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testGetWalletInfoFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newGetWalletInfoFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func newGetWalletInfoFlags(t *testing.T, testDir string) *cmd.DescribeWalletFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	return &cmd.DescribeWalletFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}
}
