package winbackoffers

import "github.com/peterbourgon/ff/v3/ffcli"

// Command returns the win-back-offers command group.
func Command() *ffcli.Command {
	return WinBackOffersCommand()
}
