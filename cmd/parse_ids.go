package cmd

import "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"

func parseCommaSeparatedIDs(input string) []string {
	return shared.SplitCSV(input)
}
