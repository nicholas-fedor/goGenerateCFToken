// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
)

// newGetCommand creates the token get subcommand.
//
// Returns:
//   - *cobra.Command: The get command that retrieves a single token's metadata.
func newGetCommand() *cobra.Command {
	gflags := &flags.GenerateFlags{}

	cmd := &cobra.Command{
		Use:   "get <token-id>",
		Short: "Get a single token's metadata",
		Long: `Retrieve full metadata for a single Cloudflare API token by ID.

Displays token name, status, issued date, expiration, and last-used timestamp.`,
		Example: `  # Get token metadata as plain text
  goGenerateCFToken token get abc123

  # Get token metadata as JSON
  goGenerateCFToken token get abc123 --json`,
		Args:    cobra.ExactArgs(1),
		GroupID: tokenGroup.ID,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCmd(cmd, args[0], gflags)
		},
	}

	gflags.BindGet(cmd.Flags())

	return cmd
}

// runGetCmd executes the token get command.
//
// Parameters:
//   - cmd: Cobra command providing context with timeout.
//   - tokenID: The token ID to retrieve.
//   - gflags: Generate flags controlling JSON output.
//
// Returns:
//   - error: Non-nil if the API key cannot be resolved or the token cannot be retrieved.
func runGetCmd(cmd *cobra.Command, tokenID string, gflags *flags.GenerateFlags) error {
	log.Debug().Str("token_id", tokenID).Msg("retrieving token metadata")

	apiKey, err := resolveAPIKey(gflags.Token)
	if err != nil {
		return err
	}

	client, err := cloudflare.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("initialize Cloudflare client: %w", err)
	}

	timeout := time.Duration(gflags.Timeout) * time.Second
	if gflags.Timeout <= 0 {
		timeout = flags.DefaultTimeout * time.Second
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
	defer cancel()

	token, err := cloudflare.GetToken(ctx, client, tokenID)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	log.Debug().Str("name", token.Name).Msg("token metadata retrieved")

	if gflags.JSON {
		data, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal JSON: %w", err)
		}

		_, werr := cmd.OutOrStdout().Write(data)
		if werr != nil {
			return fmt.Errorf("write output: %w", werr)
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		return nil
	}

	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "ID:        %s\n", token.ID)
	_, _ = fmt.Fprintf(out, "Name:      %s\n", token.Name)
	_, _ = fmt.Fprintf(out, "Status:    %s\n", token.Status)

	if !token.IssuedOn.IsZero() {
		_, _ = fmt.Fprintf(out, "IssuedOn:  %s\n", token.IssuedOn.Format(time.RFC3339))
	}

	if !token.ExpiresOn.IsZero() {
		_, _ = fmt.Fprintf(out, "ExpiresOn: %s\n", token.ExpiresOn.Format(time.RFC3339))
	}

	if !token.LastUsedOn.IsZero() {
		_, _ = fmt.Fprintf(out, "LastUsed:  %s\n", token.LastUsedOn.Format(time.RFC3339))
	}

	if !token.NotBefore.IsZero() {
		_, _ = fmt.Fprintf(out, "NotBefore: %s\n", token.NotBefore.Format(time.RFC3339))
	}

	return nil
}
