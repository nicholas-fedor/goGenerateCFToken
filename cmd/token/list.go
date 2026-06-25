// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

// newListCommand creates the token list subcommand.
//
// Returns:
//   - *cobra.Command: The list command that displays all Cloudflare API tokens.
func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List Cloudflare API tokens",
		Long: `List all Cloudflare API tokens associated with the configured credentials.

Displays token ID, name, and status for each token.`,
		Example: `  # List all tokens
  goGenerateCFToken token list`,
		GroupID: tokenGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runListCmd(cmd)
		},
	}
}

// runListCmd executes the token list command.
//
// Parameters:
//   - cmd: Cobra command providing the context with timeout.
//
// Returns:
//   - error: Non-nil if the API key cannot be resolved or tokens cannot be listed.
func runListCmd(cmd *cobra.Command) error {
	log.Debug().Msg("listing Cloudflare API tokens")

	apiKey, err := credentials.ResolveAPIKey()
	if err != nil {
		return fmt.Errorf("resolve API key: %w", err)
	}

	client, err := cloudflare.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("initialize Cloudflare client: %w", err)
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), defaultTimeout*time.Second)
	defer cancel()

	tokens, err := cloudflare.ListTokens(ctx, client)
	if err != nil {
		return fmt.Errorf("list tokens: %w", err)
	}

	log.Debug().Int("count", len(tokens)).Msg("tokens retrieved")

	for _, token := range tokens {
		logging.Printf("%s  %s  %s", token.ID, token.Name, token.Status)
	}

	return nil
}
