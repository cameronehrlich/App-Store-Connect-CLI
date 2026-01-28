package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	signingcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/signing"
)

// SigningCommand returns the signing command group.
func SigningCommand() *ffcli.Command {
	return signingcli.SigningCommand()
}
