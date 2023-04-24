// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package nodewallets_test

import (
	"context"
	"testing"

	"zuluprotocol/zeta/zeta/core/nodewallets"
	"zuluprotocol/zeta/zeta/core/nodewallets/mocks"
	vgnw "zuluprotocol/zeta/zeta/core/nodewallets/zeta"
	"zuluprotocol/zeta/zeta/core/txn"
	vgrand "zuluprotocol/zeta/zeta/libs/rand"
	vgtesting "zuluprotocol/zeta/zeta/libs/testing"
	"zuluprotocol/zeta/zeta/logging"
	commandspb "zuluprotocol/zeta/zeta/protos/zeta/commands/v1"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type testCommander struct {
	*nodewallets.Commander
	ctx    context.Context
	cfunc  context.CancelFunc
	ctrl   *gomock.Controller
	chain  *mocks.MockChain
	bstats *mocks.MockBlockchainStats
	wallet *vgnw.Wallet
}

func getTestCommander(t *testing.T) *testCommander {
	t.Helper()
	ctx, cfunc := context.WithCancel(context.Background())
	ctrl := gomock.NewController(t)
	chain := mocks.NewMockChain(ctrl)
	bstats := mocks.NewMockBlockchainStats(ctrl)
	zetaPaths, _ := vgtesting.NewZetaPaths()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	_, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass, false)
	require.NoError(t, err)
	wallet, err := nodewallets.GetZetaWallet(zetaPaths, registryPass)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	cmd, err := nodewallets.NewCommander(nodewallets.NewDefaultConfig(), logging.NewTestLogger(), chain, wallet, bstats)
	require.NoError(t, err)

	return &testCommander{
		Commander: cmd,
		ctx:       ctx,
		cfunc:     cfunc,
		ctrl:      ctrl,
		chain:     chain,
		bstats:    bstats,
		wallet:    wallet,
	}
}

func TestCommand(t *testing.T) {
	t.Run("Signed command - success", testSignedCommandSuccess)
	t.Run("Signed command - failure", testSignedCommandFailure)
}

func testSignedCommandSuccess(t *testing.T) {
	commander := getTestCommander(t)
	defer commander.Finish()

	cmd := txn.NodeVoteCommand
	payload := &commandspb.NodeVote{
		Reference: "test",
	}
	expectedChainID := vgrand.RandomStr(5)

	commander.bstats.EXPECT().Height().AnyTimes().Return(uint64(42))
	commander.chain.EXPECT().GetChainID(gomock.Any()).Times(1).Return(expectedChainID, nil)
	commander.chain.EXPECT().SubmitTransactionSync(gomock.Any(), gomock.Any()).Times(1).Return(&tmctypes.ResultBroadcastTx{}, nil)

	ok := make(chan error)
	commander.Command(context.Background(), cmd, payload, func(_ string, err error) {
		ok <- err
	}, nil)
	assert.NoError(t, <-ok)
}

func testSignedCommandFailure(t *testing.T) {
	commander := getTestCommander(t)
	defer commander.Finish()

	cmd := txn.NodeVoteCommand
	payload := &commandspb.NodeVote{
		Reference: "test",
	}
	commander.bstats.EXPECT().Height().AnyTimes().Return(uint64(42))
	commander.chain.EXPECT().GetChainID(gomock.Any()).Times(1).Return(vgrand.RandomStr(5), nil)
	commander.chain.EXPECT().SubmitTransactionSync(gomock.Any(), gomock.Any()).Times(1).Return(&tmctypes.ResultBroadcastTx{}, assert.AnError)

	ok := make(chan error)
	commander.Command(context.Background(), cmd, payload, func(_ string, err error) {
		ok <- err
	}, nil)
	assert.Error(t, <-ok)
}

func (t *testCommander) Finish() {
	t.cfunc()
	t.ctrl.Finish()
}
