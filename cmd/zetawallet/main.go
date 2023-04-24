package main

import (
	"os"

	cmd "zuluprotocol/zeta/zeta/cmd/zetawallet/commands"
)

func main() {
	writer := &cmd.Writer{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	cmd.Execute(writer)
}
