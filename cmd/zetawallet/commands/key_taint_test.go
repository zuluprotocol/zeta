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

func TestTaintKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testTaintKeyFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testTaintKeyFlagsMissingWalletFails)
	t.Run("Missing public key fails", testTaintKeyFlagsMissingPubKeyFails)
}

func testTaintKeyFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	f := &cmd.TaintKeyFlags{
		Wallet:         walletName,
		PublicKey:      pubKey,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := api.AdminTaintKeyParams{
		Wallet:     walletName,
		PublicKey:  pubKey,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedReq, req)
}

func testTaintKeyFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newTaintKeyFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func testTaintKeyFlagsMissingPubKeyFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newTaintKeyFlags(t, testDir)
	f.PublicKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("pubkey"))
	assert.Empty(t, req)
}

func newTaintKeyFlags(t *testing.T, testDir string) *cmd.TaintKeyFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.TaintKeyFlags{
		Wallet:         walletName,
		PublicKey:      pubKey,
		PassphraseFile: passphraseFilePath,
	}
}
