package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	buildbundlescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/buildbundles"
)

// BuildBundlesCommand returns the build-bundles command group.
func BuildBundlesCommand() *ffcli.Command {
	return buildbundlescli.BuildBundlesCommand()
}
