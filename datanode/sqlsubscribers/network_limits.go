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

	"github.com/pkg/errors"

	"zuluprotocol/zeta/core/events"
	"zuluprotocol/zeta/datanode/entities"
	"zuluprotocol/zeta/protos/zeta"
)

type NetworkLimitsEvent interface {
	events.Event
	NetworkLimits() *zeta.NetworkLimits
}

type NetworkLimitStore interface {
	Add(context.Context, entities.NetworkLimits) error
}

type NetworkLimits struct {
	subscriber
	store NetworkLimitStore
}

func NewNetworkLimitSub(store NetworkLimitStore) *NetworkLimits {
	t := &NetworkLimits{
		store: store,
	}
	return t
}

func (nl *NetworkLimits) Types() []events.Type {
	return []events.Type{events.NetworkLimitsEvent}
}

func (nl *NetworkLimits) Push(ctx context.Context, evt events.Event) error {
	return nl.consume(ctx, evt.(NetworkLimitsEvent))
}

func (nl *NetworkLimits) consume(ctx context.Context, event NetworkLimitsEvent) error {
	protoLimits := event.NetworkLimits()
	limits := entities.NetworkLimitsFromProto(protoLimits, entities.TxHash(event.TxHash()))
	limits.ZetaTime = nl.zetaTime

	return errors.Wrap(nl.store.Add(ctx, limits), "error adding network limits")
}
