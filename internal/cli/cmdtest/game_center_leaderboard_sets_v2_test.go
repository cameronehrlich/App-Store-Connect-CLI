package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"testing"
)

func TestGameCenterLeaderboardSetsV2ListValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "v2", "list"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
}

func TestGameCenterLeaderboardSetsV2CreateValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing app",
			args: []string{"game-center", "leaderboard-sets", "v2", "create", "--reference-name", "Season", "--vendor-id", "com.example.season"},
		},
		{
			name: "app and group",
			args: []string{"game-center", "leaderboard-sets", "v2", "create", "--app", "APP_ID", "--group-id", "GROUP_ID", "--reference-name", "Season", "--vendor-id", "com.example.season"},
		},
		{
			name: "group vendor prefix",
			args: []string{"game-center", "leaderboard-sets", "v2", "create", "--group-id", "GROUP_ID", "--reference-name", "Season", "--vendor-id", "com.example.season"},
		},
		{
			name: "missing reference-name",
			args: []string{"game-center", "leaderboard-sets", "v2", "create", "--app", "APP_ID", "--vendor-id", "com.example.season"},
		},
		{
			name: "missing vendor-id",
			args: []string{"game-center", "leaderboard-sets", "v2", "create", "--app", "APP_ID", "--reference-name", "Season"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, _ := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("expected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
		})
	}
}

func TestGameCenterLeaderboardSetMembersV2ListValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "v2", "members", "list"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
}

func TestGameCenterLeaderboardSetVersionsV2ListValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "v2", "versions", "list"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
}

func TestGameCenterLeaderboardSetLocalizationsV2CreateValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing version-id",
			args: []string{"game-center", "leaderboard-sets", "v2", "localizations", "create", "--locale", "en-US", "--name", "Season"},
		},
		{
			name: "missing locale",
			args: []string{"game-center", "leaderboard-sets", "v2", "localizations", "create", "--version-id", "VER_ID", "--name", "Season"},
		},
		{
			name: "missing name",
			args: []string{"game-center", "leaderboard-sets", "v2", "localizations", "create", "--version-id", "VER_ID", "--locale", "en-US"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, _ := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("expected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
		})
	}
}

func TestGameCenterLeaderboardSetImagesV2GetValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "v2", "images", "get"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
}

func TestGameCenterLeaderboardSetImagesV2DeleteValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing id",
			args: []string{"game-center", "leaderboard-sets", "v2", "images", "delete", "--confirm"},
		},
		{
			name: "missing confirm",
			args: []string{"game-center", "leaderboard-sets", "v2", "images", "delete", "--id", "IMAGE_ID"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, _ := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("expected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
		})
	}
}
