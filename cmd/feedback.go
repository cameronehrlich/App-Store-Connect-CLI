package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	feedbackcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/feedback"
)

// FeedbackCommand returns the feedback command group.
func FeedbackCommand() *ffcli.Command {
	return feedbackcli.FeedbackCommand()
}
