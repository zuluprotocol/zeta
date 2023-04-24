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

func TestDeleteWalletFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testDeleteWalletFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testDeleteWalletFlagsMissingWalletFails)
}

func testDeleteWalletFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	walletName := vgrand.RandomStr(10)

	f := &cmd.DeleteWalletFlags{
		Wallet: walletName,
		Force:  true,
	}

	// when
	params, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, api.AdminRemoveWalletParams{
		Wallet: walletName,
	}, params)
}

func testDeleteWalletFlagsMissingWalletFails(t *testing.T) {
	// given
	f := newDeleteWalletFlags(t)
	f.Wallet = ""

	// when
	params, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, params)
}

func newDeleteWalletFlags(t *testing.T) *cmd.DeleteWalletFlags {
	t.Helper()

	walletName := vgrand.RandomStr(10)

	return &cmd.DeleteWalletFlags{
		Wallet: walletName,
	}
}
