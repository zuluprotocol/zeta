package cmd_test

import (
	"testing"

	cmd "zuluprotocol/zeta/cmd/zetawallet/commands"
	"zuluprotocol/zeta/cmd/zetawallet/commands/flags"
	vgrand "zuluprotocol/zeta/libs/rand"
	"zuluprotocol/zeta/wallet/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListKeysFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testListKeysFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testListKeysFlagsMissingWalletFails)
}

func testListKeysFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.ListKeysFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := api.AdminListKeysParams{
		Wallet:     walletName,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedReq, req)
}

func testListKeysFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newListKeysFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func newListKeysFlags(t *testing.T, testDir string) *cmd.ListKeysFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	return &cmd.ListKeysFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}
}
