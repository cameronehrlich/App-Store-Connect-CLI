package analytics

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
)

// AnalyticsSegmentsCommand returns the analytics segments command group.
func AnalyticsSegmentsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("segments", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "segments",
		ShortUsage: "asc analytics segments <subcommand> [flags]",
		ShortHelp:  "Get analytics report segments by ID.",
		LongHelp: `Get analytics report segments by ID.

Examples:
  asc analytics segments get --segment-id "SEGMENT_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			AnalyticsSegmentsGetCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// AnalyticsSegmentsGetCommand retrieves a specific analytics report segment.
func AnalyticsSegmentsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	segmentID := fs.String("segment-id", "", "Analytics report segment ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc analytics segments get --segment-id \"SEGMENT_ID\" [flags]",
		ShortHelp:  "Get an analytics report segment by ID.",
		LongHelp: `Get an analytics report segment by ID.

Examples:
  asc analytics segments get --segment-id "SEGMENT_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if strings.TrimSpace(*segmentID) == "" {
				fmt.Fprintln(os.Stderr, "Error: --segment-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("analytics segments get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetAnalyticsReportSegment(requestCtx, strings.TrimSpace(*segmentID))
			if err != nil {
				return fmt.Errorf("analytics segments get: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}
