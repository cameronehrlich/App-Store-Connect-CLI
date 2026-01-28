package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	testflightcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/testflight"
)

// BetaTestersCommand returns the beta testers command group.
func BetaTestersCommand() *ffcli.Command {
	return testflightcli.BetaTestersCommand()
}
