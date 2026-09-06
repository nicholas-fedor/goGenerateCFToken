// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

// newRemoveCommand creates the credentials remove subcommand.
//
// Returns:
//   - *cobra.Command: The remove command that deletes the API key from the OS keyring.
func newRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove API key from OS keyring",
		Long: `Remove the stored Cloudflare API key from the OS keyring.

This does not affect any existing tokens created with the key.`,
		Example: `  # Remove the stored API key
  goGenerateCFToken credentials remove`,
		GroupID: credentialsGroup.ID,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runRemoveCmd()
		},
	}
}

// runRemoveCmd executes the credentials remove command.
//
// Returns:
//   - error: Non-nil if the API key cannot be removed from the OS keyring.
func runRemoveCmd() error {
	log.Debug().Msg("removing API key from keyring")

	err := credentials.DeleteAPIKey()
	if err != nil {
		return fmt.Errorf("remove API key: %w", err)
	}

	logging.Println("API key removed from OS keyring")

	return nil
}
