package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	prereleasecli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/prerelease"
)

// PreReleaseVersionsCommand returns the pre-release versions command group.
func PreReleaseVersionsCommand() *ffcli.Command {
	return prereleasecli.PreReleaseVersionsCommand()
}
