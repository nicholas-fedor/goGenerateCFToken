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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigFile string

var InitConfigFunc = initConfig

func InitConfig() {
	InitConfigFunc()
}

func initConfig() {
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		home := filepath.Join(homeDir, ".goGenerateCFToken")

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("CF")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("api_token", "")
	viper.SetDefault("zone", "")

	if err := viper.ReadInConfig(); err != nil {
		var configNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFoundErr) {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}
}
