/*
Copyright © 2025 Nicholas Fedor <nick@nickfedor.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nicholas-fedor/gogeneratecftoken/pkg/cloudflare"
)

var (
	// BindPFlagFunc binds a flag to a Viper key, defaulting to viper.BindPFlag.
	BindPFlagFunc = viper.BindPFlag
	// NewClientFunc creates a new Cloudflare client, defaulting to cloudflare.NewClient.
	NewClientFunc = cloudflare.NewClient
	// GenerateTokenFunc generates a Cloudflare API token, defaulting to cloudflare.GenerateToken.
	GenerateTokenFunc = cloudflare.GenerateToken
)

// generateCmd defines the command to generate a new Cloudflare API token.
var generateCmd = &cobra.Command{
	Use:   "generate [service name]",
	Short: "Generate a new Cloudflare API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		// Convert service name to lowercase for consistency.
		serviceName := strings.ToLower(args[0])

		// Retrieve API token and zone name from configuration.
		token := viper.GetString("api_token")
		zoneName := viper.GetString("zone")

		// Validate required configuration values.
		if token == "" {
			return cloudflare.ErrMissingCredentials
		}
		if zoneName == "" {
			return ErrMissingConfigZone
		}

		// Initialize Cloudflare client with the API token.
		client, err := NewClientFunc(token)
		if err != nil {
			return fmt.Errorf("failed to initialize Cloudflare client: %w", err)
		}

		// Create a context for the API call.
		ctx := context.Background()

		// Generate the new API token.
		newAPIToken, err := GenerateTokenFunc(ctx, serviceName, zoneName, client, client)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		// Output the generated token.
		fmt.Fprintln(os.Stdout, newAPIToken)

		return nil
	},
}

// init configures the generate command before execution.
func init() {
	// Add the generate command to the root command.
	rootCmd.AddCommand(generateCmd)

	// Define flags for API token and zone name.
	generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
	generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")

	// Bind the token flag to the api_token configuration key.
	if err := viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token")); err != nil {
		// Panic on binding failure, as it indicates a critical setup error.
		panic(fmt.Errorf("%w: %w", ErrBindAPITokenFlag, err))
	}

	// Bind the zone flag to the zone configuration key.
	if err := viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone")); err != nil {
		// Panic on binding failure, as it indicates a critical setup error.
		panic(fmt.Errorf("%w: %w", ErrBindZoneFlag, err))
	}
}
