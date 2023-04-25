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

package gql

import (
	"context"

	zeta "code.zetaprotocol.io/zeta/protos/zeta"
	eventspb "zuluprotocol/zeta/protos/zeta/events/v1"
)

type auctionEventResolver ZetaResolverRoot

func (r *auctionEventResolver) AuctionStart(ctx context.Context, obj *eventspb.AuctionEvent) (int64, error) {
	return obj.Start, nil
}

func (r *auctionEventResolver) AuctionEnd(ctx context.Context, obj *eventspb.AuctionEvent) (int64, error) {
	if obj.End > 0 {
		return obj.End, nil
	}
	return 0, nil
}

func (r *auctionEventResolver) ExtensionTrigger(ctx context.Context, obj *eventspb.AuctionEvent) (*zeta.AuctionTrigger, error) {
	return &obj.ExtensionTrigger, nil
}
