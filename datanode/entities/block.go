// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package entities

import (
	"encoding/hex"
	"time"

	"github.com/pkg/errors"

	"zuluprotocol/zeta/core/events"
)

type TimeUpdateEvent interface {
	events.Event
	Time() time.Time
}

type Block struct {
	ZetaTime time.Time
	Height   int64
	Hash     []byte
}

func BlockFromTimeUpdate(te TimeUpdateEvent) (*Block, error) {
	hash, err := hex.DecodeString(te.TraceID())
	if err != nil {
		return nil, errors.Wrapf(err, "Trace ID is not valid hex string, trace ID:%s", te.TraceID())
	}

	// Postgres only stores timestamps in microsecond resolution
	block := Block{
		ZetaTime: te.Time().Truncate(time.Microsecond),
		Hash:     hash,
		Height:   te.BlockNr(),
	}
	return &block, err
}
