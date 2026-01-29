package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	publishcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/publish"
)

// PublishCommand returns the publish command group.
func PublishCommand() *ffcli.Command {
	return publishcli.PublishCommand()
}
