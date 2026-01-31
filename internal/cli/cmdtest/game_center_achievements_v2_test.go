package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"testing"
)

func TestGameCenterAchievementsV2ListValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "v2", "list"}); err != nil {
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

func TestGameCenterAchievementVersionsV2ListValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "v2", "versions", "list"}); err != nil {
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

func TestGameCenterAchievementLocalizationsV2CreateValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing version-id",
			args: []string{"game-center", "achievements", "v2", "localizations", "create", "--locale", "en-US", "--name", "Test", "--before-earned-description", "Before", "--after-earned-description", "After"},
		},
		{
			name: "missing locale",
			args: []string{"game-center", "achievements", "v2", "localizations", "create", "--version-id", "VER_ID", "--name", "Test", "--before-earned-description", "Before", "--after-earned-description", "After"},
		},
		{
			name: "missing name",
			args: []string{"game-center", "achievements", "v2", "localizations", "create", "--version-id", "VER_ID", "--locale", "en-US", "--before-earned-description", "Before", "--after-earned-description", "After"},
		},
		{
			name: "missing before",
			args: []string{"game-center", "achievements", "v2", "localizations", "create", "--version-id", "VER_ID", "--locale", "en-US", "--name", "Test", "--after-earned-description", "After"},
		},
		{
			name: "missing after",
			args: []string{"game-center", "achievements", "v2", "localizations", "create", "--version-id", "VER_ID", "--locale", "en-US", "--name", "Test", "--before-earned-description", "Before"},
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

func TestGameCenterAchievementImagesV2GetValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "v2", "images", "get"}); err != nil {
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

func TestGameCenterAchievementImagesV2DeleteValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing id",
			args: []string{"game-center", "achievements", "v2", "images", "delete", "--confirm"},
		},
		{
			name: "missing confirm",
			args: []string{"game-center", "achievements", "v2", "images", "delete", "--id", "IMAGE_ID"},
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
