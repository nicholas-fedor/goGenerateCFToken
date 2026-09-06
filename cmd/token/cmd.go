// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"github.com/spf13/cobra"
)

// Token management commands.
var tokenGroup = &cobra.Group{
	ID:    "token",
	Title: "Commands:",
}

// NewCommand creates the token parent command.
//
// Returns:
//   - *cobra.Command: The token command with all subcommands attached.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Tokens",
		Long:  "Generate, list, get, and revoke Cloudflare API tokens.",
	}

	cmd.AddGroup(tokenGroup)

	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newGetCommand())
	cmd.AddCommand(newRevokeCommand())

	return cmd
}
