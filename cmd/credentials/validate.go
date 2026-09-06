// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

// newValidateCommand creates the credentials validate subcommand.
//
// Returns:
//   - *cobra.Command: The validate command that tests Cloudflare API credentials.
func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate Cloudflare credentials",
		Long: `Test the configured Cloudflare API credentials against the Cloudflare API.

Resolves the API key from CF_API_TOKEN env var or the OS keyring and verifies
it can authenticate successfully.`,
		Example: `  # Validate current credentials
  goGenerateCFToken credentials validate`,
		GroupID: credentialsGroup.ID,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runValidateCmd(cmd)
		},
	}
}

// runValidateCmd executes the credentials validate command.
//
// Parameters:
//   - cmd: Cobra command providing the context with timeout.
//
// Returns:
//   - error: Non-nil if the API key cannot be resolved or credentials are invalid.
func runValidateCmd(cmd *cobra.Command) error {
	apiKey, err := credentials.ResolveAPIKey()
	if err != nil {
		return fmt.Errorf("resolve API key: %w", err)
	}

	log.Debug().Msg("validating credentials against Cloudflare API")

	client, err := cloudflare.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("initialize Cloudflare client: %w", err)
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), flags.DefaultTimeout*time.Second)
	defer cancel()

	err = cloudflare.ValidateCredentials(ctx, client)
	if err != nil {
		return fmt.Errorf("credential validation failed: %w", err)
	}

	logging.Println("Credentials are valid")

	return nil
}
