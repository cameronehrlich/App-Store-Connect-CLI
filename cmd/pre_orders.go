package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	preorderscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/preorders"
)

// PreOrdersCommand returns the pre-orders command group.
func PreOrdersCommand() *ffcli.Command {
	return preorderscli.PreOrdersCommand()
}
