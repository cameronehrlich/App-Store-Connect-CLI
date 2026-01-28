package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	testflightcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/testflight"
)

// BetaGroupsCommand returns the beta groups command group.
func BetaGroupsCommand() *ffcli.Command {
	return testflightcli.BetaGroupsCommand()
}
