// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

// newListCommand creates the token list subcommand.
//
// Returns:
//   - *cobra.Command: The list command that displays all Cloudflare API tokens.
func newListCommand() *cobra.Command {
	var (
		outputPath   string
		jsonOutput   bool
		filterStatus string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Cloudflare API tokens",
		Long: `List all Cloudflare API tokens associated with the configured credentials.

Displays token ID, name, and status for each token.`,
		Example: `  # List all tokens as plain text
  goGenerateCFToken token list

  # List all tokens as JSON
  goGenerateCFToken token list --json

  # Filter by status
  goGenerateCFToken token list --filter active

  # Write to file
  goGenerateCFToken token list --output tokens.txt`,
		GroupID: tokenGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runListCmd(cmd, outputPath, jsonOutput, filterStatus)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Write output to file instead of stdout")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.Flags().StringVar(&filterStatus, "filter", "", "Filter tokens by status (e.g., active, disabled)")

	return cmd
}

// runListCmd executes the token list command.
//
// Parameters:
//   - cmd: Cobra command providing the context with timeout.
//   - outputPath: Optional file path for output instead of stdout.
//   - jsonOutput: If true, output as JSON.
//   - filterStatus: Optional status filter (e.g., active, disabled).
//
// Returns:
//   - error: Non-nil if the API key cannot be resolved or tokens cannot be listed.
func runListCmd(cmd *cobra.Command, outputPath string, jsonOutput bool, filterStatus string) error {
	log.Debug().Msg("listing Cloudflare API tokens")

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

	tokens, err := cloudflare.ListTokens(ctx, client)
	if err != nil {
		return fmt.Errorf("list tokens: %w", err)
	}

	log.Debug().Int("count", len(tokens)).Msg("tokens retrieved")

	if filterStatus != "" {
		tokens = filterTokensByStatus(tokens, filterStatus)
	}

	if jsonOutput {
		data, jerr := json.MarshalIndent(tokens, "", "  ")
		if jerr != nil {
			return fmt.Errorf("marshal JSON: %w", jerr)
		}

		if outputPath != "" {
			werr := os.WriteFile(outputPath, append(data, '\n'), outputFilePerms)
			if werr != nil {
				return fmt.Errorf("write output file: %w", werr)
			}

			return nil
		}

		_, werr := cmd.OutOrStdout().Write(data)
		if werr != nil {
			return fmt.Errorf("write output: %w", werr)
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		return nil
	}

	if outputPath != "" {
		var builder strings.Builder

		for _, token := range tokens {
			fmt.Fprintf(&builder, "%s  %s  %s\n", token.ID, token.Name, token.Status)
		}

		werr := os.WriteFile(outputPath, []byte(builder.String()), outputFilePerms)
		if werr != nil {
			return fmt.Errorf("write output file: %w", werr)
		}

		return nil
	}

	for _, token := range tokens {
		logging.Printf("%s  %s  %s", token.ID, token.Name, token.Status)
	}

	return nil
}

// filterTokensByStatus returns only tokens matching the given status (case-insensitive).
//
// Parameters:
//   - tokens: The full list of token info structs.
//   - status: The status string to match.
//
// Returns:
//   - []TokenInfo: The filtered list of tokens.
func filterTokensByStatus(tokens []cloudflare.TokenInfo, status string) []cloudflare.TokenInfo {
	filtered := make([]cloudflare.TokenInfo, 0, len(tokens))
	lower := strings.ToLower(status)

	for _, token := range tokens {
		if strings.ToLower(token.Status) == lower {
			filtered = append(filtered, token)
		}
	}

	return filtered
}
