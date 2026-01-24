package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

// BuildsLatestCommand returns the builds latest subcommand.
func BuildsLatestCommand() *ffcli.Command {
	fs := flag.NewFlagSet("latest", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (required, or ASC_APP_ID env)")
	version := fs.String("version", "", "Filter by version string (e.g., 1.2.3)")
	platform := fs.String("platform", "", "Filter by platform: IOS, MAC_OS, TV_OS, VISION_OS")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "latest",
		ShortUsage: "asc builds latest [flags]",
		ShortHelp:  "Get the latest build for an app.",
		LongHelp: `Get the latest build for an app.

Returns the most recently uploaded build with full metadata including
build number, version, processing state, and upload date.

This command is useful for CI/CD scripts and AI agents that need to
query the current build state before uploading a new build.

Examples:
  # Get latest build (JSON output for AI agents)
  asc builds latest --app "123456789"

  # Get latest build for a specific version
  asc builds latest --app "123456789" --version "1.2.3"

  # Filter by platform
  asc builds latest --app "123456789" --platform IOS

  # Human-readable output
  asc builds latest --app "123456789" --output table`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			// Validate platform if provided
			if strings.TrimSpace(*platform) != "" {
				validPlatforms := []string{"IOS", "MAC_OS", "TV_OS", "VISION_OS"}
				normalizedPlatform := strings.ToUpper(strings.TrimSpace(*platform))
				valid := false
				for _, p := range validPlatforms {
					if normalizedPlatform == p {
						valid = true
						break
					}
				}
				if !valid {
					fmt.Fprintf(os.Stderr, "Error: --platform must be one of: IOS, MAC_OS, TV_OS, VISION_OS\n\n")
					return flag.ErrHelp
				}
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("builds latest: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// If version is specified, we need to find the preReleaseVersion ID first
			var preReleaseVersionID string
			if strings.TrimSpace(*version) != "" {
				preReleaseVersionID, err = findPreReleaseVersionID(requestCtx, client, resolvedAppID, strings.TrimSpace(*version), strings.TrimSpace(*platform))
				if err != nil {
					return fmt.Errorf("builds latest: %w", err)
				}
				if preReleaseVersionID == "" {
					return fmt.Errorf("builds latest: no pre-release version found for version %q", *version)
				}
			}

			// Get latest build with sort by uploadedDate descending, limit 1
			opts := []asc.BuildsOption{
				asc.WithBuildsSort("-uploadedDate"),
				asc.WithBuildsLimit(1),
			}

			// Add version filter if we found a preReleaseVersion ID
			if preReleaseVersionID != "" {
				opts = append(opts, asc.WithBuildsPreReleaseVersion(preReleaseVersionID))
			}

			builds, err := client.GetBuilds(requestCtx, resolvedAppID, opts...)
			if err != nil {
				return fmt.Errorf("builds latest: failed to fetch: %w", err)
			}

			if len(builds.Data) == 0 {
				// Return empty result with appropriate message
				return fmt.Errorf("builds latest: no builds found for app %s", resolvedAppID)
			}

			// Return single build (not array) for cleaner output
			singleBuild := &asc.BuildResponse{
				Data:  builds.Data[0],
				Links: builds.Links,
			}

			return printOutput(singleBuild, *output, *pretty)
		},
	}
}

// findPreReleaseVersionID looks up the preReleaseVersion ID for a given version string.
func findPreReleaseVersionID(ctx context.Context, client *asc.Client, appID, version, platform string) (string, error) {
	opts := []asc.PreReleaseVersionsOption{
		asc.WithPreReleaseVersionsVersion(version),
		asc.WithPreReleaseVersionsLimit(1),
	}
	if platform != "" {
		opts = append(opts, asc.WithPreReleaseVersionsPlatform(platform))
	}

	versions, err := client.GetPreReleaseVersions(ctx, appID, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to lookup pre-release version: %w", err)
	}

	if len(versions.Data) == 0 {
		return "", nil
	}

	return versions.Data[0].ID, nil
}
