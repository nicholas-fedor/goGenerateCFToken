// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/prompt"
)

var (
	errFromEnvAndFile = errors.New("--from-env and --from-file are mutually exclusive")
	errEmptyKeyFile   = errors.New("API key file is empty")
	errEnvTokenUnset  = errors.New("CF_API_TOKEN is not set")
)

// newSetCommand creates the credentials set subcommand.
//
// Returns:
//   - *cobra.Command: The set command that stores an API key in the OS keyring.
func newSetCommand() *cobra.Command {
	var (
		fromEnv  bool
		fromFile string
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Store API key in the OS keyring or a local credential file",
		Long: `Store a Cloudflare API key for later commands.

The key is read interactively without echo, from --from-env (CF_API_TOKEN),
or from --from-file. It is written to the OS keyring when available, otherwise
to the default credential file under the XDG config directory.`,
		Example: `  # Store a new API key (interactive)
  goGenerateCFToken credentials set

  # Store from the CF_API_TOKEN environment variable
  goGenerateCFToken credentials set --from-env

  # Store from a file (Docker secrets)
  goGenerateCFToken credentials set --from-file /run/secrets/cf_api_token`,
		GroupID: credentialsGroup.ID,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSetCmd(fromEnv, fromFile)
		},
	}

	cmd.Flags().BoolVar(&fromEnv, "from-env", false, "Read API key from CF_API_TOKEN")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Read API key from the given file")

	return cmd
}

// runSetCmd executes the credentials set command.
//
// Parameters:
//   - fromEnv: If true, read CF_API_TOKEN instead of prompting.
//   - fromFile: If set, read the API key from this path instead of prompting.
//
// Returns:
//   - error: Non-nil if reading fails, the key is empty, or storage fails.
func runSetCmd(fromEnv bool, fromFile string) error {
	if fromEnv && fromFile != "" {
		return errFromEnvAndFile
	}

	key, err := readAPIKey(fromEnv, fromFile)
	if err != nil {
		return err
	}

	log.Debug().Msg("storing API key")

	err = credentials.SetAPIKey(key)
	if err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	logging.Println("API key stored")

	return nil
}

// readAPIKey obtains the API key from the environment, a file, or an
// interactive prompt.
//
// Parameters:
//   - fromEnv: If true, read CF_API_TOKEN instead of prompting.
//   - fromFile: If set, read the API key from this path instead of prompting.
//
// Returns:
//   - string: The API key.
//   - error: Non-nil if the source is empty or cannot be read.
func readAPIKey(fromEnv bool, fromFile string) (string, error) {
	switch {
	case fromEnv:
		key := strings.TrimSpace(os.Getenv(credentials.EnvVarToken))
		if key == "" {
			return "", errEnvTokenUnset
		}

		return key, nil
	case fromFile != "":
		data, err := os.ReadFile(fromFile)
		if err != nil {
			return "", fmt.Errorf("read API key file: %w", err)
		}

		key := strings.TrimSpace(string(data))
		if key == "" {
			return "", errEmptyKeyFile
		}

		return key, nil
	default:
		key, err := prompt.Secret("Enter Cloudflare API key: ")
		if err != nil {
			return "", fmt.Errorf("read API key: %w", err)
		}

		return key, nil
	}
}
