package cmd

import (
	"context"
	"flag"
	"os"
	"testing"
)

func TestBuildsLatestCommand_MissingApp(t *testing.T) {
	// Clear env var to ensure --app is required
	os.Unsetenv("ASC_APP_ID")

	cmd := BuildsLatestCommand()

	err := cmd.Exec(context.Background(), []string{})
	if err != flag.ErrHelp {
		t.Errorf("expected flag.ErrHelp when --app is missing, got %v", err)
	}
}

func TestBuildsLatestCommand_InvalidPlatform(t *testing.T) {
	cmd := BuildsLatestCommand()

	// Parse flags first
	if err := cmd.FlagSet.Parse([]string{"--app", "123", "--platform", "INVALID"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{})
	if err != flag.ErrHelp {
		t.Errorf("expected flag.ErrHelp for invalid platform, got %v", err)
	}
}

func TestBuildsLatestCommand_ValidPlatforms(t *testing.T) {
	validPlatforms := []string{"IOS", "MAC_OS", "TV_OS", "VISION_OS", "ios", "mac_os"}

	for _, platform := range validPlatforms {
		t.Run(platform, func(t *testing.T) {
			cmd := BuildsLatestCommand()

			// Parse flags - this should not error for valid platforms
			if err := cmd.FlagSet.Parse([]string{"--app", "123", "--platform", platform}); err != nil {
				t.Fatalf("failed to parse flags: %v", err)
			}

			// The command will fail because there's no real client, but it should get past validation
			err := cmd.Exec(context.Background(), []string{})

			// Should not be flag.ErrHelp for valid platforms (will fail later due to no auth)
			if err == flag.ErrHelp {
				t.Errorf("platform %s should be valid but got flag.ErrHelp", platform)
			}
		})
	}
}

func TestBuildsLatestCommand_FlagDefinitions(t *testing.T) {
	cmd := BuildsLatestCommand()

	// Verify all expected flags exist
	expectedFlags := []string{"app", "version", "platform", "output", "pretty"}
	for _, name := range expectedFlags {
		f := cmd.FlagSet.Lookup(name)
		if f == nil {
			t.Errorf("expected flag --%s to be defined", name)
		}
	}

	// Verify default values
	if f := cmd.FlagSet.Lookup("output"); f != nil && f.DefValue != "json" {
		t.Errorf("expected --output default to be 'json', got %q", f.DefValue)
	}
	if f := cmd.FlagSet.Lookup("pretty"); f != nil && f.DefValue != "false" {
		t.Errorf("expected --pretty default to be 'false', got %q", f.DefValue)
	}
}

func TestBuildsLatestCommand_UsesAppIDEnv(t *testing.T) {
	// Set env var
	os.Setenv("ASC_APP_ID", "env-app-id")
	defer os.Unsetenv("ASC_APP_ID")

	cmd := BuildsLatestCommand()

	// Don't pass --app flag
	if err := cmd.FlagSet.Parse([]string{}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{})

	// Should not be flag.ErrHelp since env var provides the app ID
	if err == flag.ErrHelp {
		t.Errorf("should use ASC_APP_ID env var but got flag.ErrHelp")
	}
}
