package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	routingcoveragecli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/routingcoverage"
)

// RoutingCoverageCommand returns the routing-coverage command group.
func RoutingCoverageCommand() *ffcli.Command {
	return routingcoveragecli.RoutingCoverageCommand()
}
