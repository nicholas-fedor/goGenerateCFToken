// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

var errResetForceRequired = errors.New("use --force to reset configuration")

// newResetCommand creates the config reset subcommand.
//
// Returns:
//   - *cobra.Command: The reset command that restores default configuration.
func newResetCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the config file to defaults",
		Long: `Reset the configuration file to default values.

Requires --force to confirm the destructive operation.`,
		Example: `  # Reset config to defaults
  goGenerateCFToken config reset --force`,
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResetCmd(cmd, force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Actually perform reset")

	return cmd
}

// runResetCmd executes the config reset command.
//
// Parameters:
//   - cmd: Cobra command used for flag access.
//   - force: If false, returns an error without writing.
//
// Returns:
//   - error: Non-nil if confirmation is missing or the config cannot be written.
func runResetCmd(cmd *cobra.Command, force bool) error {
	if !force {
		return errResetForceRequired
	}

	cfgPath, err := resolveConfigPath(cmd)
	if err != nil {
		return err
	}

	log.Debug().Str("path", cfgPath).Msg("resetting config file")

	err = config.EnsureFile(cfgPath)
	if err != nil {
		return fmt.Errorf("ensure config file: %w", err)
	}

	err = config.WriteDefaults(cfgPath)
	if err != nil {
		return fmt.Errorf("reset config: %w", err)
	}

	logging.Println("Configuration reset to defaults")

	return nil
}
