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

package steps

import (
	"fmt"

	"code.zetaprotocol.io/vega/core/integration/stubs"
)

func TheLiquidityFeeFactorShouldForTheMarket(
	broker *stubs.BrokerStub,
	feeStr, market string,
) error {
	mkt := broker.GetMarket(market)
	if mkt == nil {
		return fmt.Errorf("invalid market id %v", market)
	}

	got := mkt.Fees.Factors.LiquidityFee
	if got != feeStr {
		return errInvalidLiquidityFeeFactor(market, feeStr, got)
	}

	return nil
}

func errInvalidLiquidityFeeFactor(market string, expected, got string) error {
	return fmt.Errorf("invalid liquidity fee factor for market %s want %s got %s", market, expected, got)
}
