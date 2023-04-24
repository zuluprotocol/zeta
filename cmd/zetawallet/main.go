package main

import (
	"os"

	cmd "code.zetaprotocol.io/vega/cmd/vegawallet/commands"
)

func main() {
	writer := &cmd.Writer{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	cmd.Execute(writer)
}
