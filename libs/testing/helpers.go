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

package testing

import (
	"os"
	"path/filepath"

	vgrand "code.zetaprotocol.io/vega/libs/rand"
	"code.zetaprotocol.io/vega/paths"
)

func NewZetaPaths() (paths.Paths, func()) {
	path := filepath.Join("/tmp", "zeta-tests", vgrand.RandomStr(10))
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return paths.New(path), func() { _ = os.RemoveAll(path) }
}
