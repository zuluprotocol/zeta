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

package validators

import (
	"zuluprotocol/zeta/zeta/core/nodewallets"
)

type NodeWallets interface {
	GetZeta() Wallet
	GetTendermintPubkey() string
	GetEthereumAddress() string
	GetEthereum() Signer
}

type NodeWalletsWrapper struct {
	*nodewallets.NodeWallets
}

func WrapNodeWallets(nw *nodewallets.NodeWallets) *NodeWalletsWrapper {
	return &NodeWalletsWrapper{nw}
}

func (w *NodeWalletsWrapper) GetZeta() Wallet {
	return w.Zeta
}

func (w *NodeWalletsWrapper) GetEthereum() Signer {
	return w.Ethereum
}

func (w *NodeWalletsWrapper) GetEthereumAddress() string {
	return w.Ethereum.PubKey().Hex()
}

func (w *NodeWalletsWrapper) GetTendermintPubkey() string {
	return w.Tendermint.Pubkey
}
