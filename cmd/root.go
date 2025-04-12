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
	"os"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/config"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "goGenerateCFToken",
	Version: version,
	Short:   "Cloudflare API token management tool",
	Long: `goGenerateCFToken is a CLI tool for managing Cloudflare API tokens.
It currently supports generating API tokens scoped for Zone Read/DNS Write.

Usage instructions:
1) Create a config.yaml configuration file:
	Example:
		api_token: "your-cloudflare-api-token-here"
		zone: "example.com"
2) Run the following command:
	goGenerateCFToken generate [your service]
3) If successful, the API token will be printed to the console
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVar(
		&config.ConfigFile,
		"config",
		"",
		"config file (default is $HOME/.goGenerateCFToken/config.yaml)",
	)
}
