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
	"code.zetaprotocol.io/vega/logging"
)

func DebugOrders(broker *stubs.BrokerStub, log *logging.Logger) {
	log.Info("DUMPING ORDERS")
	data := broker.GetOrderEvents()
	for _, v := range data {
		o := *v.Order()
		log.Info(fmt.Sprintf("order %s: %v\n", o.Id, o.String()))
	}
}