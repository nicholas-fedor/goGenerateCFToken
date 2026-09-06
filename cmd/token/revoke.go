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
)

var errForceRequired = errors.New("use --force to revoke token")

// newRevokeCommand creates the token revoke subcommand.
//
// Returns:
//   - *cobra.Command: The revoke command that deletes a Cloudflare API token by ID.
func newRevokeCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "revoke <token-id>",
		Short: "Revoke a Cloudflare API token",
		Long: `Revoke (delete) a Cloudflare API token by its ID.

Requires --force to confirm the destructive operation.`,
		Example: `  # Revoke a token
  goGenerateCFToken token revoke abc123 --force`,
		Args:    cobra.ExactArgs(1),
		GroupID: tokenGroup.ID,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRevokeCmd(cmd, args, force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Confirm token revocation")

	return cmd
}

// runRevokeCmd executes the token revoke command.
//
// Parameters:
//   - cmd: Cobra command providing the context with timeout.
//   - args: Positional arguments. args[0] is the token ID to revoke.
//   - force: If false, logs a warning and returns without revoking.
//
// Returns:
//   - error: Non-nil if the API key cannot be resolved or revocation fails.
func runRevokeCmd(cmd *cobra.Command, args []string, force bool) error {
	if !force {
		return errForceRequired
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
