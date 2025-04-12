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

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/cloudflare"
)

var (
	NewClientFunc     = cloudflare.NewClient
	GenerateTokenFunc = (*cloudflare.Client).GenerateToken
)

var generateCmd = &cobra.Command{
	Use:   "generate [service name]",
	Short: "Generate a new Cloudflare API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		serviceName := strings.ToLower(args[0])
		token := viper.GetString("api_token")
		zoneName := viper.GetString("zone")

		if token == "" {
			return ErrMissingConfigAuth
		}
		if zoneName == "" {
			return ErrMissingConfigZone
		}

		client, err := NewClientFunc(token)
		if err != nil {
			return fmt.Errorf("failed to initialize Cloudflare client: %w", err)
		}

		ctx := context.Background()

		newAPIToken, err := GenerateTokenFunc(client, ctx, serviceName, zoneName)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		fmt.Fprintln(os.Stdout, newAPIToken)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
	generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")

	if err := viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token")); err != nil {
		panic(fmt.Errorf("%w: %w", ErrBindAPITokenFlag, err))
	}

	if err := viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone")); err != nil {
		panic(fmt.Errorf("%w: %w", ErrBindZoneFlag, err))
	}
}
