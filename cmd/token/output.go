// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/flags"
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
	asJSON := gflags.JSON || gflags.Format == "json"
	writeFile := gflags.Output != ""
	writeStdout := gflags.Format != "none" && !writeFile

	if !writeFile && !writeStdout {
		return nil
	}

	if asJSON {
		return writeTokenJSON(cmd, gflags, result, writeFile, writeStdout)
	}

	return writeTokenPlain(cmd, gflags, result, writeFile, writeStdout)
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
	return writeTokenJSON(cmd, &flags.GenerateFlags{}, result, false, true)
}

// marshalTokenJSON serializes the token generation result as indented JSON.
//
// Parameters:
//   - result: The token generation result to serialize.
//
// Returns:
//   - []byte: Indented JSON with a trailing newline.
//   - error: Non-nil if JSON marshaling fails.
func marshalTokenJSON(result *cloudflare.TokenGenerationResult) ([]byte, error) {
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
		return nil, fmt.Errorf("marshal JSON: %w", err)
	}

	return append(jsonData, '\n'), nil
}

// writeTokenJSON writes the token result as indented JSON to a file and/or
// stdout.
//
// Parameters:
//   - cmd: Cobra command used for stdout output.
//   - gflags: Generate flags providing the optional output file path.
//   - result: The token generation result to serialize.
//   - writeFile: If true, write JSON to gflags.Output.
//   - writeStdout: If true, write JSON to stdout.
//
// Returns:
//   - error: Non-nil if JSON marshaling or writing fails.
func writeTokenJSON(
	cmd *cobra.Command,
	gflags *flags.GenerateFlags,
	result *cloudflare.TokenGenerationResult,
	writeFile, writeStdout bool,
) error {
	jsonData, err := marshalTokenJSON(result)
	if err != nil {
		return err
	}

	if writeFile {
		log.Debug().Str("path", gflags.Output).Msg("writing token JSON to file")

		werr := os.WriteFile(gflags.Output, jsonData, outputFilePerms)
		if werr != nil {
			return fmt.Errorf("write output file: %w", werr)
		}
	}

	if writeStdout {
		_, werr := cmd.OutOrStdout().Write(jsonData)
		if werr != nil {
			return fmt.Errorf("write output: %w", werr)
		}
	}

	return nil
}

// writeTokenPlain writes the raw token value to a file and/or stdout.
//
// Parameters:
//   - cmd: Cobra command used for stdout output.
//   - gflags: Generate flags providing the optional output file path.
//   - result: The token generation result containing the token value.
//   - writeFile: If true, write the token to gflags.Output.
//   - writeStdout: If true, write the token to stdout.
//
// Returns:
//   - error: Non-nil if writing fails.
func writeTokenPlain(
	cmd *cobra.Command,
	gflags *flags.GenerateFlags,
	result *cloudflare.TokenGenerationResult,
	writeFile, writeStdout bool,
) error {
	body := []byte(result.Token + "\n")

	if writeFile {
		log.Debug().Str("path", gflags.Output).Msg("writing token to file")

		werr := os.WriteFile(gflags.Output, body, outputFilePerms)
		if werr != nil {
			return fmt.Errorf("write output file: %w", werr)
		}
	}

	if writeStdout {
		_, werr := cmd.OutOrStdout().Write(body)
		if werr != nil {
			return fmt.Errorf("write output: %w", werr)
		}
	}

	return nil
}
