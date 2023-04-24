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

package sqlsubscribers

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"zuluprotocol/zeta/zeta/core/events"
	"zuluprotocol/zeta/zeta/datanode/entities"
	types "zuluprotocol/zeta/zeta/protos/zeta"
)

type TradeEvent interface {
	events.Event
	Trade() types.Trade
}

type TradesStore interface {
	Add(*entities.Trade) error
	Flush(ctx context.Context) error
}

type TradeSubscriber struct {
	subscriber
	store TradesStore
}

func NewTradesSubscriber(store TradesStore) *TradeSubscriber {
	return &TradeSubscriber{
		store: store,
	}
}

func (ts *TradeSubscriber) Types() []events.Type {
	return []events.Type{events.TradeEvent}
}

func (ts *TradeSubscriber) Flush(ctx context.Context) error {
	return ts.store.Flush(ctx)
}

func (ts *TradeSubscriber) Push(ctx context.Context, evt events.Event) error {
	return ts.consume(evt.(TradeEvent))
}

func (ts *TradeSubscriber) consume(te TradeEvent) error {
	trade := te.Trade()
	return errors.Wrap(ts.addTrade(&trade, entities.TxHash(te.TxHash()), ts.zetaTime, te.Sequence()), "failed to consume trade")
}

func (ts *TradeSubscriber) addTrade(t *types.Trade, txHash entities.TxHash, zetaTime time.Time, blockSeqNumber uint64) error {
	trade, err := entities.TradeFromProto(t, txHash, zetaTime, blockSeqNumber)
	if err != nil {
		return errors.Wrap(err, "converting event to trade")
	}

	return errors.Wrap(ts.store.Add(trade), "adding trade to store")
}
