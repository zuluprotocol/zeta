package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"zuluprotocol/zeta/cmd/zetawallet/commands/cli"
	"zuluprotocol/zeta/cmd/zetawallet/commands/flags"
	"zuluprotocol/zeta/cmd/zetawallet/commands/printer"
	"zuluprotocol/zeta/wallet/api"

	"github.com/spf13/cobra"
)

var (
	verifyMessageLong = cli.LongDesc(`
		Verify a message against a signature.

		The signature has to be generated by Ed25519 cryptographic algorithm.
	`)

	verifyMessageExample = cli.Examples(`
		# Verify the signature of a message
		{{.Software}} message verify --message MESSAGE --signature SIGNATURE --pubkey PUBKEY
	`)
)

type VerifyMessageHandler func(api.AdminVerifyMessageParams) (api.AdminVerifyMessageResult, error)

func NewCmdVerifyMessage(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(params api.AdminVerifyMessageParams) (api.AdminVerifyMessageResult, error) {
		verifyMessage := api.NewAdminVerifyMessage()
		rawResult, errorDetails := verifyMessage.Handle(context.Background(), params)
		if errorDetails != nil {
			return api.AdminVerifyMessageResult{}, errors.New(errorDetails.Data)
		}
		return rawResult.(api.AdminVerifyMessageResult), nil
	}
	return BuildCmdVerifyMessage(w, h, rf)
}

func BuildCmdVerifyMessage(w io.Writer, handler VerifyMessageHandler, rf *RootFlags) *cobra.Command {
	f := &VerifyMessageFlags{}

	cmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verify a message against a signature",
		Long:    verifyMessageLong,
		Example: verifyMessageExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			resp, err := handler(req)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintVerifyMessageResponse(w, resp.IsValid)
			case flags.JSONOutput:
				return printer.FprintJSON(w, struct {
					IsValid bool `json:"isValid"`
				}{
					IsValid: resp.IsValid,
				})
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.PubKey,
		"pubkey", "k",
		"",
		"Public key associated to the signature (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.Message,
		"message", "m",
		"",
		"Message to be verified (base64-encoded)",
	)
	cmd.Flags().StringVarP(&f.Signature,
		"signature", "s",
		"",
		"Signature of the message (base64-encoded)",
	)

	return cmd
}

type VerifyMessageFlags struct {
	Signature string
	Message   string
	PubKey    string
}

func (f *VerifyMessageFlags) Validate() (api.AdminVerifyMessageParams, error) {
	req := api.AdminVerifyMessageParams{}

	if len(f.PubKey) == 0 {
		return api.AdminVerifyMessageParams{}, flags.MustBeSpecifiedError("pubkey")
	}
	req.PubKey = f.PubKey

	if len(f.Signature) == 0 {
		return api.AdminVerifyMessageParams{}, flags.MustBeSpecifiedError("signature")
	}
	_, err := base64.StdEncoding.DecodeString(f.Signature)
	if err != nil {
		return api.AdminVerifyMessageParams{}, flags.MustBase64EncodedError("signature")
	}
	req.EncodedSignature = f.Signature

	if len(f.Message) == 0 {
		return api.AdminVerifyMessageParams{}, flags.MustBeSpecifiedError("message")
	}
	_, err = base64.StdEncoding.DecodeString(f.Message)
	if err != nil {
		return api.AdminVerifyMessageParams{}, flags.MustBase64EncodedError("message")
	}
	req.EncodedMessage = f.Message

	return req, nil
}

func PrintVerifyMessageResponse(w io.Writer, isValid bool) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	if isValid {
		str.CheckMark().SuccessText("Valid signature").NextSection()
	} else {
		str.CrossMark().DangerText("Invalid signature").NextSection()
	}

	str.BlueArrow().InfoText("Sign a message").NextLine()
	str.Text("To sign a message, see the following command:").NextSection()
	str.Code(fmt.Sprintf("%s sign --help", os.Args[0])).NextLine()
}
