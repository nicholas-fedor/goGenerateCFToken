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
	"strings"

	"github.com/nicholas-fedor/goGenerateCFToken/cloudflare"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	NewAPIClientFunc  = cloudflare.NewAPIClient
	GenerateTokenFunc = cloudflare.GenerateToken
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:   "generate [service name]",
	Short: "Generate a new Cloudflare API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := strings.ToLower(args[0])
		token := viper.GetString("api_token")
		zone := viper.GetString("zone")

		if token == "" {
			return fmt.Errorf("missing API token: set via --token, CF_API_TOKEN, or config file")
		}
		if zone == "" {
			return fmt.Errorf("missing zone: set via --zone, CF_ZONE, or config file")
		}

		// Create API client using the scoped API token
		client, err := NewAPIClientFunc(token)
		if err != nil {
			return fmt.Errorf("failed to create API client: %v", err)
		}

		// Most API calls require a Context
		ctx := context.Background()

		// Generate an API token
		newAPIToken, err := GenerateTokenFunc(serviceName, zone, client, ctx)
		if err != nil {
			return fmt.Errorf("failed to generate token: %v", err)
		}

		fmt.Println(newAPIToken)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
	generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")
	viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token"))
	viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone"))
}
