// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"github.com/spf13/cobra"
)

// Credential management commands.
var credentialsGroup = &cobra.Group{
	ID:    "credentials",
	Title: "Commands:",
}

// NewCommand creates the credentials parent command.
//
// Returns:
//   - *cobra.Command: The credentials command with all subcommands attached.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Credentials",
		Long:  "Store, remove, and validate Cloudflare API credentials.",
	}

	cmd.AddGroup(credentialsGroup)

	cmd.AddCommand(newSetCommand())
	cmd.AddCommand(newRemoveCommand())
	cmd.AddCommand(newValidateCommand())

	return cmd
}
