package assets_test

import (
	"testing"

	"zuluprotocol/zeta/core/assets"
	erc20mocks "zuluprotocol/zeta/core/assets/erc20/mocks"
	"zuluprotocol/zeta/core/assets/mocks"
	bmocks "zuluprotocol/zeta/core/broker/mocks"
	"zuluprotocol/zeta/core/nodewallets"
	nweth "zuluprotocol/zeta/core/nodewallets/eth"
	nwzeta "code.zetaprotocol.io/zeta/core/nodewallets/zeta"
	"zuluprotocol/zeta/core/types"
	"zuluprotocol/zeta/libs/num"
	vgrand "zuluprotocol/zeta/libs/rand"
	"zuluprotocol/zeta/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testService struct {
	*assets.Service
	broker     *bmocks.MockInterface
	bridgeView *mocks.MockERC20BridgeView
	notary     *mocks.MockNotary
	ctrl       *gomock.Controller
}

func TestAssets(t *testing.T) {
	t.Run("Staging asset update for unknown asset fails", testStagingAssetUpdateForUnknownAssetFails)
}

func testStagingAssetUpdateForUnknownAssetFails(t *testing.T) {
	service := getTestService(t)

	// given
	asset := &types.Asset{
		ID: vgrand.RandomStr(5),
		Details: &types.AssetDetails{
			Name:     vgrand.RandomStr(5),
			Symbol:   vgrand.RandomStr(3),
			Decimals: 10,
			Quantum:  num.DecimalFromInt64(42),
			Source: &types.AssetDetailsErc20{
				ERC20: &types.ERC20{
					ContractAddress:   vgrand.RandomStr(5),
					LifetimeLimit:     num.NewUint(42),
					WithdrawThreshold: num.NewUint(84),
				},
			},
		},
	}

	// when
	err := service.StageAssetUpdate(asset)

	// then
	require.ErrorIs(t, err, assets.ErrAssetDoesNotExist)
}

func getTestService(t *testing.T) *testService {
	t.Helper()
	conf := assets.NewDefaultConfig()
	logger := logging.NewTestLogger()
	ctrl := gomock.NewController(t)
	ethClient := erc20mocks.NewMockETHClient(ctrl)
	broker := bmocks.NewMockInterface(ctrl)
	bridgeView := mocks.NewMockERC20BridgeView(ctrl)
	notary := mocks.NewMockNotary(ctrl)
	nodeWallets := &nodewallets.NodeWallets{
		Zeta:     &nwzeta.Wallet{},
		Ethereum: &nweth.Wallet{},
	}
	service := assets.New(logger, conf, nodeWallets, ethClient, broker, bridgeView, notary, true)
	return &testService{
		Service:    service,
		broker:     broker,
		ctrl:       ctrl,
		bridgeView: bridgeView,
		notary:     notary,
	}
}
