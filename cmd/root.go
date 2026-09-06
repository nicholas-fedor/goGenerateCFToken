// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package cmd is the root CLI command.
package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/cmd/config"
	"github.com/nicholas-fedor/gogeneratecftoken/cmd/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/cmd/token"
	"github.com/nicholas-fedor/gogeneratecftoken/cmd/version"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

var rootCmd = &cobra.Command{
	Use:   "goGenerateCFToken",
	Short: "A CLI for Cloudflare API token management",
	Long: `goGenerateCFToken creates and manages Cloudflare API tokens with DNS edit permissions.

The configuration is loaded from an XDG-compliant config file
($XDG_CONFIG_HOME/gogeneratecftoken/config.yaml) or flags.

The Cloudflare API token is resolved from CF_API_TOKEN, CF_API_TOKEN_FILE, the OS
keyring, or the default credential file.`,
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		return setupLogging(cmd)
	}
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	cflags := &flags.CommonFlags{}
	cflags.Bind(rootCmd.PersistentFlags())

	rootCmd.AddCommand(config.NewCommand())
	rootCmd.AddCommand(credentials.NewCommand())
	rootCmd.AddCommand(token.NewCommand())
	rootCmd.AddCommand(token.NewDeprecatedGenerateCommand())
	rootCmd.AddCommand(version.NewCommand())
}

// Execute initializes and runs the root CLI command.
//
// Returns:
//   - error: Non-nil if command execution fails.
func Execute() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	rootCmd.SetContext(ctx)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	return nil
}

// Root exposes the root command for external tools like doc generators.
//
// Returns:
//   - *cobra.Command: The root command instance.
func Root() *cobra.Command {
	return rootCmd
}

// setupLogging configures zerolog based on command flags.
//
// Parameters:
//   - cmd: Cobra command providing flag values.
//
// Returns:
//   - error: Non-nil if flag retrieval fails.
func setupLogging(cmd *cobra.Command) error {
	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		return fmt.Errorf("get quiet flag: %w", err)
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return fmt.Errorf("get verbose flag: %w", err)
	}

	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return fmt.Errorf("get log-level flag: %w", err)
	}

	switch {
	case quiet:
		logging.SetQuiet(true)
		logging.Setup(zerolog.ErrorLevel)
	case verbose:
		logging.SetQuiet(false)
		logging.Setup(zerolog.DebugLevel)
	default:
		level, perr := zerolog.ParseLevel(logLevel)
		if perr != nil {
			return fmt.Errorf("parse log level %q: %w", logLevel, perr)
		}

		logging.SetQuiet(false)
		logging.Setup(level)
	}

	return nil
}
