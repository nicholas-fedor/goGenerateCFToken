// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

// newSetCommand creates the credentials set subcommand.
//
// Returns:
//   - *cobra.Command: The set command that stores an API key in the OS keyring.
func newSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Store API key in OS keyring",
		Long: `Prompt for a Cloudflare API key and store it securely in the OS keyring.

The key is used by other commands to authenticate with the Cloudflare API.`,
		Example: `  # Store a new API key
  goGenerateCFToken credentials set`,
		GroupID: credentialsGroup.ID,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSetCmd()
		},
	}
}

// runSetCmd executes the credentials set command, prompting for an API key.
//
// Returns:
//   - error: Non-nil if reading from stdin fails, key is empty, or the API key cannot be stored.
func runSetCmd() error {
	fmt.Fprint(os.Stderr, "Enter Cloudflare API key: ")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read API key: %w", err)
	}

	key := strings.TrimSpace(input)
	if key == "" {
		return errors.New("API key cannot be empty")
	}

	log.Debug().Msg("storing API key in keyring")

	err = credentials.SetAPIKey(key)
	if err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	logging.Println("API key stored in OS keyring")

	return nil
}
