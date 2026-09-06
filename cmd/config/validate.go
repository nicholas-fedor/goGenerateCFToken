// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/logging"
)

// newValidateCommand creates the config validate subcommand.
//
// Returns:
//   - *cobra.Command: The validate command that checks the config file for errors.
func newValidateCommand() *cobra.Command {
	cflags := &flags.ConfigFlags{}

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration file",
		Long: `Check the configuration file for errors.

Verifies the file exists, is valid YAML, and contains a valid zone name.`,
		Example: `  # Validate the default config file
  goGenerateCFToken config validate

  # Validate a specific config file
  goGenerateCFToken config validate --config ./my-config.yaml`,
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runConfigValidateCmd(cmd, cflags)
		},
	}

	cflags.Bind(cmd.Flags())

	return cmd
}

// runConfigValidateCmd executes the config validate command.
//
// Parameters:
//   - cmd: Cobra command used for flag access.
//   - cflags: Config flags providing the output format.
//
// Returns:
//   - error: Non-nil if the config file cannot be loaded or is invalid.
func runConfigValidateCmd(cmd *cobra.Command, _ *flags.ConfigFlags) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("get config flag: %w", err)
	}

	log.Debug().Str("config_path", configPath).Msg("validating configuration")

	_, err = config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logging.Println("Configuration is valid")

	return nil
}
