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

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	vgjson "code.zetaprotocol.io/vega/libs/json"
	"code.zetaprotocol.io/vega/version"
)

const (
	outputFlagName     = "output"
	outputFlagValJSON  = "json"
	outputFlagValHuman = "human"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().String(outputFlagName, outputFlagValHuman, "Specify the output format: json,human")
}

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Returns a Zeta Visor version",
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := cmd.Flags().GetString(outputFlagName)
		if err != nil {
			return err
		}

		switch output {
		case outputFlagValHuman:
			fmt.Printf("Zeta Visor CLI %s (%s)\n", version.Get(), version.GetCommitHash())
			return nil
		case outputFlagValJSON:
			return vgjson.Print(struct {
				Version string `json:"version"`
				Hash    string `json:"hash"`
			}{
				Version: version.Get(),
				Hash:    version.GetCommitHash(),
			})
		default:
			return fmt.Errorf("%s flag must be either %q or %q", outputFlagName, outputFlagValHuman, outputFlagValJSON)
		}
	},
}
