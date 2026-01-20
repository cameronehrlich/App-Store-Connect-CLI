package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/auth"
)

// Auth command factory
func AuthCommand() *ffcli.Command {
	fs := flag.NewFlagSet("auth", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "auth",
		ShortUsage: "asc auth <subcommand> [flags]",
		ShortHelp:  "Manage App Store Connect API authentication.",
		LongHelp:   "Manage App Store Connect API authentication.\n\nAuthentication is handled via App Store Connect API keys. Generate keys at:\nhttps://appstoreconnect.apple.com/access/integrations/api\n\nCredentials are stored in the system keychain when available, with a local config fallback.\n\nSubcommands:\n  login     Register and store API key\n  logout    Remove stored credentials\n  status    Show current authentication status",
		FlagSet:    fs,
		Subcommands: []*ffcli.Command{
			AuthLoginCommand(),
			AuthLogoutCommand(),
			AuthStatusCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			if len(args) == 0 {
				fs.Usage()
				return flag.ErrHelp
			}
			return nil
		},
	}
}

// AuthLogin command factory
func AuthLoginCommand() *ffcli.Command {
	fs := flag.NewFlagSet("auth login", flag.ExitOnError)

	name := fs.String("name", "", "Friendly name for this key")
	keyID := fs.String("key-id", "", "App Store Connect API Key ID")
	issuerID := fs.String("issuer-id", "", "App Store Connect Issuer ID")
	keyPath := fs.String("private-key", "", "Path to private key (.p8) file")

	return &ffcli.Command{
		Name:       "login",
		ShortUsage: "asc auth login [flags]",
		ShortHelp:  "Register and store App Store Connect API key.",
		LongHelp: `Register and store App Store Connect API key.

This command stores your API credentials in the system keychain when available,
with a local config fallback (restricted permissions).

Examples:
  asc auth login --name "MyKey" --key-id "ABC123" --issuer-id "DEF456" --private-key /path/to/AuthKey.p8

The private key file path is stored securely. The key content is never saved.`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := fs.Parse(args); err != nil {
				return err
			}

			if *name == "" {
				fmt.Fprintln(os.Stderr, "Error: --name is required")
				fs.Usage()
				return flag.ErrHelp
			}
			if *keyID == "" {
				fmt.Fprintln(os.Stderr, "Error: --key-id is required")
				fs.Usage()
				return flag.ErrHelp
			}
			if *issuerID == "" {
				fmt.Fprintln(os.Stderr, "Error: --issuer-id is required")
				fs.Usage()
				return flag.ErrHelp
			}
			if *keyPath == "" {
				fmt.Fprintln(os.Stderr, "Error: --private-key is required")
				fs.Usage()
				return flag.ErrHelp
			}

			// Validate the key file exists and is parseable
			if err := auth.ValidateKeyFile(*keyPath); err != nil {
				return fmt.Errorf("invalid private key: %w", err)
			}

			// Store credentials securely
			if err := auth.StoreCredentials(*name, *keyID, *issuerID, *keyPath); err != nil {
				return fmt.Errorf("failed to store credentials: %w", err)
			}

			fmt.Printf("Successfully registered API key '%s'\n", *name)
			return nil
		},
	}
}

// AuthLogout command factory
func AuthLogoutCommand() *ffcli.Command {
	fs := flag.NewFlagSet("auth logout", flag.ExitOnError)
	all := fs.Bool("all", false, "Remove all stored credentials (default)")

	return &ffcli.Command{
		Name:       "logout",
		ShortUsage: "asc auth logout [flags]",
		ShortHelp:  "Remove stored API credentials.",
		LongHelp: `Remove stored API credentials.

Examples:
  asc auth logout
  asc auth logout --all`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := fs.Parse(args); err != nil {
				return err
			}

			if *all {
				// Flag is accepted for future multi-key support.
			}

			if err := auth.RemoveAllCredentials(); err != nil {
				return fmt.Errorf("failed to remove credentials: %w", err)
			}

			fmt.Println("Successfully removed stored credentials")
			return nil
		},
	}
}

// AuthStatus command factory
func AuthStatusCommand() *ffcli.Command {
	fs := flag.NewFlagSet("auth status", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "status",
		ShortUsage: "asc auth status",
		ShortHelp:  "Show current authentication status.",
		LongHelp: `Show current authentication status.

Displays information about stored API keys and which one is currently active.

Examples:
  asc auth status`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := fs.Parse(args); err != nil {
				return err
			}

			credentials, err := auth.ListCredentials()
			if err != nil {
				return fmt.Errorf("failed to list credentials: %w", err)
			}

			if len(credentials) == 0 {
				fmt.Println("No credentials stored. Run 'asc auth login' to get started.")
				return nil
			}

			fmt.Println("Stored credentials:")
			for _, cred := range credentials {
				active := ""
				if cred.IsDefault {
					active = " (default)"
				}
				fmt.Printf("  - %s (Key ID: %s)%s\n", cred.Name, cred.KeyID, active)
			}
			return nil
		},
	}
}
