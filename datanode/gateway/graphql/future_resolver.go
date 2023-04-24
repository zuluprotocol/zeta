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

	protoTypes "code.zetaprotocol.io/vega/protos/vega"
	zetapb "code.vegaprotocol.io/vega/protos/vega"
)

type myFutureResolver ZetaResolverRoot

func (r *myFutureResolver) SettlementAsset(ctx context.Context, obj *protoTypes.Future) (*protoTypes.Asset, error) {
	return r.r.getAssetByID(ctx, obj.SettlementAsset)
}

func (r *myFutureResolver) DataSourceSpecForSettlementData(_ context.Context, obj *zetapb.Future) (*DataSourceSpec, error) {
	if obj.DataSourceSpecForSettlementData == nil {
		return nil, nil
	}
	return resolveDataSourceSpec(obj.DataSourceSpecForSettlementData), nil
}

func (r *myFutureResolver) DataSourceSpecForTradingTermination(_ context.Context, obj *zetapb.Future) (*DataSourceSpec, error) {
	if obj.DataSourceSpecForTradingTermination == nil {
		return nil, nil
	}
	return resolveDataSourceSpec(obj.DataSourceSpecForTradingTermination), nil
}
