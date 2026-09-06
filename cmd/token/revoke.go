// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/prompt"
)

var errRevokeCancelled = errors.New("revocation cancelled")

// newRevokeCommand creates the token revoke subcommand.
//
// Returns:
//   - *cobra.Command: The revoke command that deletes a Cloudflare API token by ID.
func newRevokeCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "revoke <token-id>",
		Short: "Revoke a Cloudflare API token",
		Long: `Revoke (delete) a Cloudflare API token by its ID.

Prompts for confirmation (y/N) unless --yes or -y is provided.`,
		Example: `  # Revoke a token (with confirmation prompt)
  goGenerateCFToken token revoke abc123

  # Revoke a token without confirmation
  goGenerateCFToken token revoke abc123 --yes`,
		Args:    cobra.ExactArgs(1),
		GroupID: tokenGroup.ID,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRevokeCmd(cmd, args, yes)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Bypass confirmation prompt")

	return cmd
}

// runRevokeCmd executes the token revoke command.
//
// Parameters:
//   - cmd: Cobra command providing the context with timeout.
//   - args: Positional arguments. args[0] is the token ID to revoke.
//   - yes: If true, bypass the confirmation prompt.
//
// Returns:
//   - error: Non-nil if the user declines, the API key cannot be resolved, or revocation fails.
func runRevokeCmd(cmd *cobra.Command, args []string, yes bool) error {
	if !yes {
		ok, err := prompt.Confirm("Revoke token " + args[0] + "? (y/N) ")
		if err != nil || !ok {
			return errRevokeCancelled
		}
	}

	apiKey, err := credentials.ResolveAPIKey()
	if err != nil {
		return fmt.Errorf("resolve API key: %w", err)
	}

	client, err := cloudflare.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("initialize Cloudflare client: %w", err)
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), flags.DefaultTimeout*time.Second)
	defer cancel()

	err = cloudflare.RevokeToken(ctx, client, args[0])
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}

	logging.Printf("Token revoked: %s", args[0])

	return nil
}
