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

	"zuluprotocol/zeta/zeta/core/integration/stubs"
	types "zuluprotocol/zeta/zeta/protos/zeta"
)

func PartyShouldHaveOneAccountPerAsset(
	broker *stubs.BrokerStub,
	owner string,
) error {
	assets := map[string]struct{}{}

	accounts := broker.GetAccounts()

	for _, acc := range accounts {
		if acc.Owner == owner && acc.Type == types.AccountType_ACCOUNT_TYPE_GENERAL {
			if _, ok := assets[acc.Asset]; ok {
				return errMultipleGeneralAccountForAsset(owner, acc)
			}
			assets[acc.Asset] = struct{}{}
		}
	}
	return nil
}

func errMultipleGeneralAccountForAsset(owner string, acc types.Account) error {
	return fmt.Errorf("party=%v have multiple account for asset=%v", owner, acc.Asset)
}
