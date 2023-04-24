package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/cli"
	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/flags"
	"zuluprotocol/zeta/zeta/cmd/zetawallet/commands/printer"
	"zuluprotocol/zeta/zeta/wallet/api"
	"zuluprotocol/zeta/zeta/wallet/wallets"

	"github.com/spf13/cobra"
)

var (
	listWalletsLong = cli.LongDesc(`
		List all registered wallets.
	`)

	listWalletsExample = cli.Examples(`
		# List all registered wallets
		{{.Software}} list
	`)
)

type ListWalletsHandler func() (api.AdminListWalletsResult, error)

func NewCmdListWallets(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func() (api.AdminListWalletsResult, error) {
		walletStore, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return api.AdminListWalletsResult{}, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}
		defer walletStore.Close()

		listWallet := api.NewAdminListWallets(walletStore)
		rawResult, errorDetails := listWallet.Handle(context.Background(), nil)
		if errorDetails != nil {
			return api.AdminListWalletsResult{}, errors.New(errorDetails.Data)
		}
		return rawResult.(api.AdminListWalletsResult), nil
	}

	return BuildCmdListWallets(w, h, rf)
}

func BuildCmdListWallets(w io.Writer, handler ListWalletsHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all registered wallets",
		Long:    listWalletsLong,
		Example: listWalletsExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			resp, err := handler()
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintListWalletsResult(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	return cmd
}

func PrintListWalletsResult(w io.Writer, resp api.AdminListWalletsResult) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	if len(resp.Wallets) == 0 {
		str.InfoText("No wallet registered").NextLine()
		return
	}

	for _, w := range resp.Wallets {
		str.Text(fmt.Sprintf("- %s", w)).NextLine()
	}
}
