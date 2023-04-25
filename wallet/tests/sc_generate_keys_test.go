package tests_test

import (
	"testing"

	vgrand "zuluprotocol/zeta/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestGenerateKey(t *testing.T) {
	// given
	home := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, home)
	walletName := vgrand.RandomStr(5)

	// when
	createWalletResp, err := WalletCreate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertCreateWallet(t, createWalletResp).
		WithName(walletName).
		LocatedUnder(home)

	// when
	descResp, err := KeyDescribe(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", createWalletResp.Key.PublicKey,
	})

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{"name": "Key 1"}).
		WithAlgorithm("zeta/ed25519", 1).
		WithTainted(false)

	// when
	generateKeyResp, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:key-2,role:validation",
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithMetadata(map[string]string{"name": "key-2", "role": "validation"})

	// when
	descResp, err = KeyDescribe(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", generateKeyResp.PublicKey,
	})

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{"name": "key-2", "role": "validation"}).
		WithAlgorithm("zeta/ed25519", 1).
		WithTainted(false)
}
