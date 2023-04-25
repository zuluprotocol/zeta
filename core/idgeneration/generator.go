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

package idgeneration

import (
	"encoding/hex"

	"zuluprotocol/zeta/libs/crypto"
)

// IDGenerator no mutex required, markets work deterministically, and sequentially.
type IDGenerator struct {
	nextIDBytes []byte
}

// New returns an idGenerator, and is used to abstract this type.
func New(rootID string) *IDGenerator { //revive:disable:unexported-return
	nextIDBytes, err := hex.DecodeString(rootID)
	if err != nil {
		panic("failed to create new deterministic id generator: " + err.Error())
	}

	return &IDGenerator{
		nextIDBytes: nextIDBytes,
	}
}

func (i *IDGenerator) NextID() string {
	if i == nil {
		panic("id generator instance is not initialised")
	}

	nextID := hex.EncodeToString(i.nextIDBytes)
	i.nextIDBytes = crypto.Hash(i.nextIDBytes)
	return nextID
}
