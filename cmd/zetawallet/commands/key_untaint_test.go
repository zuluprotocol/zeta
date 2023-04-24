package cmd_test

import (
	"testing"

	cmd "code.zetaprotocol.io/vega/cmd/vegawallet/commands"
	"code.zetaprotocol.io/vega/cmd/vegawallet/commands/flags"
	vgrand "code.zetaprotocol.io/vega/libs/rand"
	"code.zetaprotocol.io/vega/wallet/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUntaintKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testUntaintKeyFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testUntaintKeyFlagsMissingWalletFails)
	t.Run("Missing public key fails", testUntaintKeyFlagsMissingPubKeyFails)
}

func testUntaintKeyFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	f := &cmd.UntaintKeyFlags{
		Wallet:         walletName,
		PublicKey:      pubKey,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := api.AdminUntaintKeyParams{
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

func testUntaintKeyFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newUntaintKeyFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func testUntaintKeyFlagsMissingPubKeyFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newUntaintKeyFlags(t, testDir)
	f.PublicKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("pubkey"))
	assert.Empty(t, req)
}

func newUntaintKeyFlags(t *testing.T, testDir string) *cmd.UntaintKeyFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.UntaintKeyFlags{
		Wallet:         walletName,
		PublicKey:      pubKey,
		PassphraseFile: passphraseFilePath,
	}
}
