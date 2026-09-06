// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"github.com/spf13/cobra"
)

// Configuration management commands.
var configGroup = &cobra.Group{
	ID:    "config",
	Title: "Commands:",
}

// NewCommand creates the config parent command.
//
// Returns:
//   - *cobra.Command: The config command with all subcommands attached.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration",
		Long:  "View, set, reset, delete, and validate the configuration file.",
	}

	cmd.AddGroup(configGroup)

	cmd.AddCommand(newSetCommand())
	cmd.AddCommand(newShowCommand())
	cmd.AddCommand(newDeleteCommand())
	cmd.AddCommand(newResetCommand())
	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newInitCommand())

	return cmd
}
