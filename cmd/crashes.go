package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	crashescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/crashes"
)

// CrashesCommand returns the crashes command group.
func CrashesCommand() *ffcli.Command {
	return crashescli.CrashesCommand()
}
