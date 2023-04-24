// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package nodewallet

import (
	"fmt"

	"code.zetaprotocol.io/vega/core/config"
	"code.zetaprotocol.io/vega/core/nodewallets"
	"code.zetaprotocol.io/vega/logging"
	"code.zetaprotocol.io/vega/paths"

	"github.com/jessevdk/go-flags"
)

type verifyCmd struct {
	Config nodewallets.Config
}

func (opts *verifyCmd) Execute(_ []string) error {
	log := logging.NewLoggerFromConfig(logging.NewDefaultConfig())
	defer log.AtExit()

	registryPass, err := rootCmd.PassphraseFile.Get("node wallet", false)
	if err != nil {
		return err
	}

	zetaPaths := paths.New(rootCmd.ZetaHome)

	_, conf, err := config.EnsureNodeConfig(zetaPaths)
	if err != nil {
		return err
	}

	opts.Config = conf.NodeWallet

	if _, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse(); err != nil {
		return err
	}

	nw, err := nodewallets.GetNodeWallets(opts.Config, zetaPaths, registryPass)
	if err != nil {
		return fmt.Errorf("couldn't get node wallets: %w", err)
	}

	if err := nw.Verify(); err != nil {
		return err
	}

	fmt.Println(green("ok"))
	return nil
}
