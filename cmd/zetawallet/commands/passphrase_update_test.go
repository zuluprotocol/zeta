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

func TestUpdatePassphrase(t *testing.T) {
	t.Run("Valid flags succeeds", testUpdatePassphraseFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testUpdatePassphraseFlagsMissingWalletFails)
}

func testUpdatePassphraseFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	isolatedPassphrase, isolatedPassphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.UpdatePassphraseFlags{
		Wallet:            walletName,
		PassphraseFile:    passphraseFilePath,
		NewPassphraseFile: isolatedPassphraseFilePath,
	}

	expectedReq := api.AdminUpdatePassphraseParams{
		Wallet:        walletName,
		Passphrase:    passphrase,
		NewPassphrase: isolatedPassphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedReq, req)
}

func testUpdatePassphraseFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newUpdatePassphraseFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func newUpdatePassphraseFlags(t *testing.T, testDir string) *cmd.UpdatePassphraseFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	_, newPassphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	return &cmd.UpdatePassphraseFlags{
		Wallet:            walletName,
		NewPassphraseFile: newPassphraseFilePath,
		PassphraseFile:    passphraseFilePath,
	}
}
