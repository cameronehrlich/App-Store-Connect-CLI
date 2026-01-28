package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	testflightcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/testflight"
)

// TestFlightCommand returns the testflight command group.
func TestFlightCommand() *ffcli.Command {
	return testflightcli.TestFlightCommand()
}
