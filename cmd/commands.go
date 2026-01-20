package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/auth"
)

// Feedback command factory
func FeedbackCommand() *ffcli.Command {
	fs := flag.NewFlagSet("feedback", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	jsonFlag := fs.Bool("json", false, "Output in JSON format (shorthand)")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "feedback",
		ShortUsage: "asc feedback [flags]",
		ShortHelp:  "List TestFlight feedback from beta testers.",
		LongHelp: `List TestFlight feedback from beta testers.

This command fetches beta feedback screenshot submissions and comments.

Examples:
  asc feedback --app "123456789"
  asc feedback --app "123456789" --json`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := fs.Parse(args); err != nil {
				return err
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			feedback, err := client.GetFeedback(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("failed to fetch feedback: %w", err)
			}

			format := *output
			if *jsonFlag {
				format = "json"
			}

			return printOutput(feedback, format, *pretty)
		},
	}
}

// Crashes command factory
func CrashesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("crashes", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	jsonFlag := fs.Bool("json", false, "Output in JSON format (shorthand)")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "crashes",
		ShortUsage: "asc crashes [flags]",
		ShortHelp:  "List and export TestFlight crash reports.",
		LongHelp: `List and export TestFlight crash reports.

This command fetches crash reports submitted by TestFlight beta testers,
helping you identify and fix issues in your app.

Examples:
  asc crashes --app "123456789"
  asc crashes --app "123456789" --json
  asc crashes --app "123456789" --json > crashes.json`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := fs.Parse(args); err != nil {
				return err
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			crashes, err := client.GetCrashes(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("failed to fetch crashes: %w", err)
			}

			format := *output
			if *jsonFlag {
				format = "json"
			}

			return printOutput(crashes, format, *pretty)
		},
	}
}

// Reviews command factory
func ReviewsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("reviews", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	jsonFlag := fs.Bool("json", false, "Output in JSON format (shorthand)")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")
	stars := fs.Int("stars", 0, "Filter by star rating (1-5)")
	territory := fs.String("territory", "", "Filter by territory (e.g., US, GBR)")

	return &ffcli.Command{
		Name:       "reviews",
		ShortUsage: "asc reviews [flags]",
		ShortHelp:  "List App Store customer reviews.",
		LongHelp: `List App Store customer reviews.

This command fetches customer reviews from the App Store,
helping you understand user feedback and sentiment.

Examples:
  asc reviews --app "123456789"
  asc reviews --app "123456789" --json
  asc reviews --app "123456789" --stars 1 --territory US --json`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := fs.Parse(args); err != nil {
				return err
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.ReviewOption{}
			if *stars != 0 {
				if *stars < 1 || *stars > 5 {
					return fmt.Errorf("--stars must be between 1 and 5")
				}
				opts = append(opts, asc.WithRating(*stars))
			}
			if *territory != "" {
				opts = append(opts, asc.WithTerritory(*territory))
			}

			reviews, err := client.GetReviews(requestCtx, resolvedAppID, opts...)
			if err != nil {
				return fmt.Errorf("failed to fetch reviews: %w", err)
			}

			format := *output
			if *jsonFlag {
				format = "json"
			}

			return printOutput(reviews, format, *pretty)
		},
	}
}

// RootCommand returns the root command
func RootCommand(version string) *ffcli.Command {
	root := &ffcli.Command{
		Name:       "asc",
		ShortUsage: "asc <subcommand> [flags]",
		ShortHelp:  "A fast, AI-agent friendly CLI for App Store Connect.",
		LongHelp:   "ASC is a lightweight CLI for App Store Connect. Built for developers and AI agents.",
		FlagSet:    flag.NewFlagSet("asc", flag.ExitOnError),
		Subcommands: []*ffcli.Command{
			AuthCommand(),
			FeedbackCommand(),
			CrashesCommand(),
			ReviewsCommand(),
		},
	}

	versionFlag := root.FlagSet.Bool("version", false, "Print version and exit")

	root.Exec = func(ctx context.Context, args []string) error {
		if *versionFlag {
			fmt.Fprintln(root.FlagSet.Output(), version)
			return nil
		}
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		}
		return flag.ErrHelp
	}

	return root
}

func getASCClient() (*asc.Client, error) {
	actualKeyID := os.Getenv("ASC_KEY_ID")
	actualIssuerID := os.Getenv("ASC_ISSUER_ID")
	actualKeyPath := os.Getenv("ASC_PRIVATE_KEY_PATH")

	if actualKeyID == "" || actualIssuerID == "" || actualKeyPath == "" {
		cfg, err := auth.GetDefaultCredentials()
		if err != nil && actualKeyID == "" && actualIssuerID == "" && actualKeyPath == "" {
			return nil, err
		}
		if cfg != nil {
			if actualKeyID == "" {
				actualKeyID = cfg.KeyID
			}
			if actualIssuerID == "" {
				actualIssuerID = cfg.IssuerID
			}
			if actualKeyPath == "" {
				actualKeyPath = cfg.PrivateKeyPath
			}
		}
	}

	if actualKeyID == "" || actualIssuerID == "" || actualKeyPath == "" {
		return nil, fmt.Errorf("missing authentication. Run 'asc auth login'")
	}

	return asc.NewClient(actualKeyID, actualIssuerID, actualKeyPath)
}

func printOutput(data interface{}, format string, pretty bool) error {
	format = strings.ToLower(format)
	switch format {
	case "json":
		if pretty {
			return asc.PrintPrettyJSON(data)
		}
		return asc.PrintJSON(data)
	case "markdown", "md":
		if pretty {
			return fmt.Errorf("--pretty is only valid with JSON output")
		}
		return asc.PrintMarkdown(data)
	case "table":
		if pretty {
			return fmt.Errorf("--pretty is only valid with JSON output")
		}
		return asc.PrintTable(data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func resolveAppID(appID string) string {
	if appID != "" {
		return appID
	}
	return os.Getenv("ASC_APP_ID")
}

func contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, asc.DefaultTimeout)
}
