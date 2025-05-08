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
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/pkg/config"
)

var (
	// goos holds the operating system type for determining configuration paths.
	// It defaults to runtime.GOOS but can be overridden for testing.
	goos = runtime.GOOS
	// configFilePath specifies the default configuration file location.
	configFilePath = configPath()
	// shortDescription provides a brief summary of the CLI tool.
	shortDescription = "A CLI generator for Cloudflare API Tokens"
	// longDescription provides detailed usage information for the CLI tool.
	longDescription = fmt.Sprintf(
		`
			goGenerateCFToken

A CLI tool for creating Cloudflare API tokens with DNS edit permissions.

Configuration Filepath:
  %s

Example Configuration
  api_token: your-cloudflare-api-token-here
  zone: example.com

Instructions:
1) Create a configuration file.
2) Run the following command:
     goGenerateCFToken generate [desired prefix]
3) The API token will be created and printed to the console

Example:
  Zone: example.com
  Command: goGenerateCFToken generate service
  Output:
	Token Name: service.example.com
	Token Value: (random token value)

`, configFilePath)
)

// rootCmd defines the root command for the CLI tool.
var rootCmd = &cobra.Command{
	Use:   "goGenerateCFToken",
	Short: shortDescription,
	Long:  longDescription,
}

// Execute runs the root command, handling errors by exiting with a non-zero status.
func Execute() {
	// Execute the root command and check for errors.
	if err := rootCmd.Execute(); err != nil {
		// Exit with status 1 on error.
		os.Exit(1)
	}
}

// SetVersionInfo sets the version information for the root command.
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

// init configures the root command before execution.
func init() {
	// Initialize configuration loading on command start.
	cobra.OnInitialize(config.InitConfig)

	// Define the persistent --config flag for specifying the configuration file.
	rootCmd.PersistentFlags().StringVar(
		&config.ConfigFile,
		config.ConfigFilename,
		"",
		configFilePath,
	)
}

// userHomeDir returns the user’s home directory based on the operating system.
// It returns an empty string for unsupported operating systems.
func userHomeDir() string {
	switch goos {
	case "windows":
		return "%userprofile%"
	case "linux":
		return "$HOME"
	case "ios":
		return "/"
	case "plan9":
		return "$home"
	case "android":
		return "/sdcard"
	}
	// Return empty string for unknown OS.
	return ""
}

// configPath returns the default configuration file path for the CLI tool.
func configPath() string {
	// Construct the path using the home directory, app directory, and config file name.
	return fmt.Sprintf(
		"%s.%s",
		filepath.Join(userHomeDir(), config.AppDirName, config.ConfigFilename),
		config.ConfigExt,
	)
}
