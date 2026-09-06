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
	"github.com/nicholas-fedor/gogeneratecftoken/internal/prompt"
)

var errResetCancelled = errors.New("reset cancelled")

// newResetCommand creates the config reset subcommand.
//
// Returns:
//   - *cobra.Command: The reset command that restores default configuration.
func newResetCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the config file to defaults",
		Long: `Reset the configuration file to default values.

Prompts for confirmation (y/N) unless --yes or -y is provided.`,
		Example: `  # Reset config (with confirmation prompt)
  goGenerateCFToken config reset

  # Reset config without confirmation
  goGenerateCFToken config reset --yes`,
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResetCmd(cmd, yes)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Bypass confirmation prompt")

	return cmd
}

// runResetCmd executes the config reset command.
//
// Parameters:
//   - cmd: Cobra command used for flag access.
//   - yes: If true, bypass the confirmation prompt.
//
// Returns:
//   - error: Non-nil if the user declines or the config cannot be written.
func runResetCmd(cmd *cobra.Command, yes bool) error {
	if !yes {
		ok, err := prompt.Confirm("Reset configuration to defaults? (y/N) ")
		if err != nil || !ok {
			return errResetCancelled
		}
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
