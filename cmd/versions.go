package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	versionscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/versions"
)

// VersionsCommand returns the versions command group.
func VersionsCommand() *ffcli.Command {
	return versionscli.VersionsCommand()
}
