// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
)

// outputFilePerms defines secure file permissions for written token files (owner read/write only).
const outputFilePerms = 0o600

// outputToken writes the generated token to stdout or a file.
//
// Parameters:
//   - cmd: Cobra command used for output.
//   - gflags: Generate flags controlling JSON output and optional file path.
//   - result: The token generation result containing the token value and metadata.
//
// Returns:
//   - error: Non-nil if writing to stdout or the output file fails.
func outputToken(cmd *cobra.Command, gflags *flags.GenerateFlags, result *cloudflare.TokenGenerationResult) error {
	if gflags.JSON {
		return outputTokenJSON(cmd, result)
	}

	return outputTokenPlain(cmd, gflags, result)
}

// outputTokenJSON writes the token result as indented JSON to stdout.
//
// Parameters:
//   - cmd: Cobra command used for stdout output.
//   - result: The token generation result to serialize.
//
// Returns:
//   - error: Non-nil if JSON marshaling or stdout write fails.
func outputTokenJSON(cmd *cobra.Command, result *cloudflare.TokenGenerationResult) error {
	data := map[string]string{
		"token": result.Token,
		"name":  result.Name,
		"zone":  result.Zone,
	}

	if result.ExpiresOn != "" {
		data["expiresOn"] = result.ExpiresOn
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	_, werr := cmd.OutOrStdout().Write(jsonData)
	if werr != nil {
		return fmt.Errorf("write output: %w", werr)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	return nil
}

// outputTokenPlain writes the token as plain text to stdout or a file.
//
// Parameters:
//   - cmd: Cobra command used for stdout output.
//   - gflags: Generate flags providing the optional output file path.
//   - result: The token generation result containing the token value.
//
// Returns:
//   - error: Non-nil if writing to stdout or the output file fails.
func outputTokenPlain(cmd *cobra.Command, gflags *flags.GenerateFlags, result *cloudflare.TokenGenerationResult) error {
	if gflags.Output != "" {
		log.Debug().Str("path", gflags.Output).Msg("writing token to file")

		werr := os.WriteFile(gflags.Output, []byte(result.Token+"\n"), outputFilePerms)
		if werr != nil {
			return fmt.Errorf("write output file: %w", werr)
		}
	} else {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), result.Token)
	}

	return nil
}
