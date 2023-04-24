package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"code.zetaprotocol.io/vega/cmd/vegawallet/commands/flags"
	"code.zetaprotocol.io/vega/cmd/vegawallet/commands/printer"
	vgterm "code.zetaprotocol.io/vega/libs/term"
	vgzap "code.zetaprotocol.io/vega/libs/zap"
	"code.zetaprotocol.io/vega/paths"
	"code.zetaprotocol.io/vega/wallet/api"
	netstore "code.zetaprotocol.io/vega/wallet/network/store/v1"
	"code.zetaprotocol.io/vega/wallet/wallets"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	DefaultForwarderRetryCount = 5
)

type Error struct {
	Err string `json:"error"`
}

type Writer struct {
	Out io.Writer
	Err io.Writer
}

func Execute(w *Writer) {
	c := NewCmdRoot(w.Out)

	execErr := c.Execute()
	if execErr == nil {
		return
	}

	defer os.Exit(1)

	if errors.Is(execErr, flags.ErrUnsupportedOutput) {
		_, _ = fmt.Fprintln(w.Err, execErr)
	}

	output, _ := c.Flags().GetString("output")
	switch output {
	case flags.InteractiveOutput:
		fprintErrorInteractive(w, execErr)
	case flags.JSONOutput:
		fprintErrorJSON(w.Err, execErr)
	}
}

func fprintErrorInteractive(w *Writer, execErr error) {
	if vgterm.HasTTY() {
		p := printer.NewInteractivePrinter(w.Out)
		p.Print(p.String().CrossMark().DangerText("Error: ").DangerText(execErr.Error()).NextLine())
	} else {
		_, _ = fmt.Fprintln(w.Err, execErr)
	}
}

func fprintErrorJSON(w io.Writer, err error) {
	jsonErr := printer.FprintJSON(w, Error{
		Err: err.Error(),
	})
	if jsonErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "couldn't format error as JSON: %v\n", jsonErr)
		_, _ = fmt.Fprintf(os.Stderr, "original error: %v\n", err)
	}
}

func autoCompleteWallet(cmd *cobra.Command, zetaHome string, property string) {
	err := cmd.RegisterFlagCompletionFunc(property, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		walletStore, err := wallets.InitialiseStore(zetaHome)
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		defer walletStore.Close()

		listWallet := api.NewAdminListWallets(walletStore)
		rawResult, errorDetails := listWallet.Handle(context.Background(), nil)
		if errorDetails != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return rawResult.(api.AdminListWalletsResult).Wallets, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		panic(err)
	}
}

func autoCompleteNetwork(cmd *cobra.Command, zetaHome string) {
	err := cmd.RegisterFlagCompletionFunc("network", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		zetaPaths := paths.New(vegaHome)

		netStore, err := netstore.InitialiseStore(zetaPaths)
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}

		nets, err := netStore.ListNetworks()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return nets, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		panic(err)
	}
}

func autoCompleteLogLevel(cmd *cobra.Command) {
	err := cmd.RegisterFlagCompletionFunc("level", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return vgzap.SupportedLogLevels, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		panic(err)
	}
}

func buildCmdLogger(output, level string) (*zap.Logger, error) {
	if output == flags.InteractiveOutput {
		return vgzap.BuildStandardConsoleLogger(level)
	}

	return vgzap.BuildStandardJSONLogger(level)
}