package commands_test

import (
	"encoding/hex"
	"testing"

	"zuluprotocol/zeta/commands"
	commandspb "zuluprotocol/zeta/protos/zeta/commands/v1"
	"github.com/stretchr/testify/assert"
)

func TestCheckAnnounceNode(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilAnnounceNodeFails)
	t.Run("Submitting a node registration without zeta pub key fails", testAnnounceNodeWithoutZetaPubKeyFails)
	t.Run("Submitting a node registration with valid zeta pub key succeeds", testAnnounceNodeWithValidZetaPubKeySucceeds)
	t.Run("Submitting a node registration with invalid encoding of zeta pub key succeeds", testAnnounceNodeWithInvalidEncodingOfZetaPubKeyFails)
	t.Run("Submitting a node registration without ethereum pub key fails", testAnnounceNodeWithoutEthereumAddressFails)
	t.Run("Submitting a node registration with ethereum address succeeds", testAnnounceNodeWithEthereumAddressSucceeds)
	t.Run("Submitting a node registration without chain address fails", testAnnounceNodeWithoutChainPubKeyFails)
	t.Run("Submitting a node registration with chain pub key succeeds", testAnnounceNodeWithChainPubKeySucceeds)
}

func testNilAnnounceNodeFails(t *testing.T) {
	err := checkAnnounceNode(nil)

	assert.Error(t, err)
}

func testAnnounceNodeWithoutZetaPubKeyFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.zeta_pub_key"), commands.ErrIsRequired)
}

func testAnnounceNodeWithValidZetaPubKeySucceeds(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		ZetaPubKey: hex.EncodeToString([]byte("0xDEADBEEF")),
	})
	assert.NotContains(t, err.Get("announce_node.zeta_pub_key"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("announce_node.zeta_pub_key"), commands.ErrShouldBeHexEncoded)
}

func testAnnounceNodeWithInvalidEncodingOfZetaPubKeyFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		ZetaPubKey: "invalid-hex-encoding",
	})
	assert.Contains(t, err.Get("announce_node.zeta_pub_key"), commands.ErrShouldBeAValidZetaPubkey)
}

func testAnnounceNodeWithoutEthereumAddressFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.ethereum_address"), commands.ErrIsRequired)
}

func testAnnounceNodeWithEthereumAddressSucceeds(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		EthereumAddress: "0xDEADBEEF",
	})
	assert.NotContains(t, err.Get("announce_node.ethereum_address"), commands.ErrIsRequired)
}

func testAnnounceNodeWithoutChainPubKeyFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.chain_pub_key"), commands.ErrIsRequired)
}

func testAnnounceNodeWithChainPubKeySucceeds(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		ChainPubKey: "0xDEADBEEF",
	})
	assert.NotContains(t, err.Get("announce_node.chain_pub_key"), commands.ErrIsRequired)
}

func checkAnnounceNode(cmd *commandspb.AnnounceNode) commands.Errors {
	err := commands.CheckAnnounceNode(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
