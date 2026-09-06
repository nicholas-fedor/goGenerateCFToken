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

var errDeletionCancelled = errors.New("deletion cancelled")

// newDeleteCommand creates the config delete subcommand.
//
// Returns:
//   - *cobra.Command: The delete command that removes the config file and directory.
func newDeleteCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the config file and directory",
		Long: `Delete the configuration file and its parent directory.

Prompts for confirmation (y/N) unless --yes or -y is provided.`,
		Example: `  # Delete config (with confirmation prompt)
  goGenerateCFToken config delete

  # Delete config without confirmation
  goGenerateCFToken config delete --yes`,
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDeleteCmd(cmd, yes)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Bypass confirmation prompt")

	return cmd
}

// runDeleteCmd executes the config delete command.
//
// Parameters:
//   - cmd: Cobra command used for flag access.
//   - yes: If true, bypass the confirmation prompt.
//
// Returns:
//   - error: Non-nil if the user declines or the config cannot be deleted.
func runDeleteCmd(cmd *cobra.Command, yes bool) error {
	if !yes {
		ok, err := prompt.Confirm("Delete config file and directory? (y/N) ")
		if err != nil || !ok {
			return errDeletionCancelled
		}
	}

	cfgPath, err := resolveConfigPath(cmd)
	if err != nil {
		return err
	}

	log.Debug().Str("path", cfgPath).Msg("deleting config file")

	err = config.Delete(cfgPath)
	if err != nil {
		return fmt.Errorf("delete config: %w", err)
	}

	logging.Println("Configuration deleted")

	return nil
}
