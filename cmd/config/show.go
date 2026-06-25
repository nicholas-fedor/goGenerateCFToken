// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
)

// newShowCommand creates the config show subcommand.
//
// Returns:
//   - *cobra.Command: The show command that displays current configuration values.
func newShowCommand() *cobra.Command {
	cflags := &flags.ConfigFlags{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long: `Display the current configuration values.

Supports YAML (default) and JSON output formats.`,
		Example: `  # Show config as YAML
  goGenerateCFToken config show

  # Show config as JSON
  goGenerateCFToken config show --format json`,
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runShowCmd(cmd, cflags)
		},
	}

	cflags.Bind(cmd.Flags())
	cflags.BindShow(cmd.Flags())

	return cmd
}

// runShowCmd executes the config show command.
//
// Parameters:
//   - cmd: Cobra command used for output and flag access.
//   - cflags: Config flags providing the output format.
//
// Returns:
//   - error: Non-nil if the config cannot be loaded or formatted.
func runShowCmd(cmd *cobra.Command, cflags *flags.ConfigFlags) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("get config flag: %w", err)
	}

	log.Debug().Str("format", cflags.Format).Msg("loading configuration")

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	data, err := cfg.Format(cflags.Format)
	if err != nil {
		return fmt.Errorf("format config: %w", err)
	}

	cmd.Print(string(data))

	return nil
}
