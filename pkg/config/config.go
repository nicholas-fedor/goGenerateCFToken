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

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Constants defining the configuration directory and file details.
const (
	// AppDirName is the directory name for storing configuration files.
	AppDirName = ".goGenerateCFToken"
	// ConfigFilename is the base name of the configuration file.
	ConfigFilename = "config"
	// ConfigExt is the file extension for the configuration file.
	ConfigExt = "yaml"
)

var (
	// ConfigFile specifies the path to the configuration file, set via flag or default.
	ConfigFile string

	// InitConfigFunc is the function to initialize configuration, injectable for testing.
	InitConfigFunc func(Viper)

	// osUserHomeDir retrieves the user’s home directory, defaulting to os.UserHomeDir.
	osUserHomeDir = os.UserHomeDir

	// osExit terminates the program with an exit code, defaulting to os.Exit.
	osExit = os.Exit
)

// Viper defines the interface for configuration management, wrapping viper.Viper methods.
type Viper interface {
	// SetConfigFile sets the explicit configuration file path.
	SetConfigFile(file string)

	// AddConfigPath adds a directory to search for configuration files.
	AddConfigPath(path string)

	// SetConfigName sets the base name of the configuration file.
	SetConfigName(name string)

	// SetConfigType sets the configuration file format (e.g., yaml).
	SetConfigType(typ string)

	// SetEnvPrefix sets the prefix for environment variables.
	SetEnvPrefix(prefix string)

	// SetEnvKeyReplacer sets the replacer for environment variable keys.
	SetEnvKeyReplacer(r *strings.Replacer)

	// AutomaticEnv enables automatic binding of environment variables.
	AutomaticEnv()

	// SetDefault sets a default value for a configuration key.
	SetDefault(key string, value any)

	// ReadInConfig loads the configuration from the specified file.
	ReadInConfig() error

	// ConfigFileUsed returns the path of the loaded configuration file.
	ConfigFileUsed() string
}

// InitConfig initializes the configuration using Viper.
func InitConfig() {
	// Retrieve the Viper instance.
	v := viper.GetViper()

	// Set default initialization function if none provided.
	if InitConfigFunc == nil {
		InitConfigFunc = initConfig
	}

	// Execute the initialization function.
	InitConfigFunc(v)
}

// initConfig configures Viper with file paths, environment variables, and loads the config.
func initConfig(v Viper) {
	// Set up configuration file paths and name.
	if err := setConfigFile(v); err != nil {
		// Report error and exit if file setup fails.
		fmt.Fprintf(os.Stderr, "Failed to set config file: %v\n", err)
		osExit(1)
	}

	// Configure environment variable bindings.
	setEnv(v)

	// Load the configuration file.
	loadConfig(v)
}

// setConfigFile configures Viper with the configuration file path or default locations.
// It returns an error if the home directory cannot be determined.
func setConfigFile(v Viper) error {
	switch {
	case ConfigFile != "":
		// Use explicit configuration file if specified.
		v.SetConfigFile(ConfigFile)

	default:
		// Determine user’s home directory for default config path.
		homeDir, err := osUserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		// Construct default config path in home directory.
		configPath := filepath.Join(homeDir, AppDirName)

		// Add config paths: home directory and current directory.
		v.AddConfigPath(configPath)
		v.AddConfigPath(".")

		// Set config file name and type.
		v.SetConfigName(ConfigFilename)
		v.SetConfigType(ConfigExt)
	}

	return nil
}

// setEnv configures Viper to bind environment variables with defaults.
func setEnv(v Viper) {
	// Set environment variable prefix to "CF".
	v.SetEnvPrefix("CF")

	// Replace dots with underscores in environment variable keys.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable automatic environment variable binding.
	v.AutomaticEnv()

	// Set default values for required configuration keys.
	v.SetDefault("api_token", "")
	v.SetDefault("zone", "")
}

// loadConfig attempts to load the configuration file and reports its status.
func loadConfig(v Viper) {
	// Try to read the configuration file.
	if err := v.ReadInConfig(); err != nil {
		// Handle non-"file not found" errors by reporting to stderr.
		var configNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFoundErr) {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		}
	} else {
		// Report the loaded configuration file path.
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", v.ConfigFileUsed())
	}
}
