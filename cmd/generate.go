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
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nicholas-fedor/goGenerateCFToken/cloudflare"
)

var (
	NewAPIClientFunc  = cloudflare.NewAPIClient
	GenerateTokenFunc = cloudflare.GenerateToken
)

// Static error variables.
var (
	ErrMissingAPIToken = errors.New(
		"missing API token: set via --token, CF_API_TOKEN, or config file",
	)
	ErrMissingZone         = errors.New("missing zone: set via --zone, CF_ZONE, or config file")
	ErrCreateClientFailed  = errors.New("failed to create API client")
	ErrGenerateTokenFailed = errors.New("failed to generate token")
	ErrBindAPITokenFlag    = errors.New("failed to bind api_token flag")
	ErrBindZoneFlag        = errors.New("failed to bind zone flag")
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:   "generate [service name]",
	Short: "Generate a new Cloudflare API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		serviceName := strings.ToLower(args[0])
		token := viper.GetString("api_token")
		zone := viper.GetString("zone")

		if token == "" {
			return ErrMissingAPIToken
		}
		if zone == "" {
			return ErrMissingZone
		}

		// Create API client using the scoped API token
		client, err := NewAPIClientFunc(token)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreateClientFailed, err)
		}

		// Most API calls require a Context
		ctx := context.Background()

		// Generate an API token
		newAPIToken, err := GenerateTokenFunc(ctx, serviceName, zone, client)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGenerateTokenFailed, err)
		}

		fmt.Println(newAPIToken) //nolint:forbidigo // Intended output behavior

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
	generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")

	// Bind flags to viper and check for errors
	if err := viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token")); err != nil {
		panic(fmt.Errorf("%w: %w", ErrBindAPITokenFlag, err))
	}

	if err := viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone")); err != nil {
		panic(fmt.Errorf("%w: %w", ErrBindZoneFlag, err))
	}
}
