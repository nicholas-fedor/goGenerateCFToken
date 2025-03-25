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
	"testing"

	"github.com/nicholas-fedor/goGenerateCFToken/config"
)

func TestRootCmd(t *testing.T) {
	// Test that rootCmd is initialized correctly
	if rootCmd.Use != "goGenerateCFToken" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "goGenerateCFToken")
	}

	if rootCmd.Version != version {
		t.Errorf("rootCmd.Version = %q, want %q", rootCmd.Version, version)
	}

	if rootCmd.Short == "" || rootCmd.Long == "" {
		t.Errorf("rootCmd Short or Long description is empty")
	}

	// Test Execute with mocked os.Exit
	oldExit := osExit
	osExit = func(code int) {
		if code != 0 {
			t.Errorf("Execute() exited with code %d", code)
		}
	}

	defer func() { osExit = oldExit }()

	// Set a dummy config file to avoid real config loading
	oldConfigFile := config.ConfigFile
	config.ConfigFile = "dummy.yaml"

	defer func() { config.ConfigFile = oldConfigFile }()

	Execute()

	// Test persistent flags
	configFileFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFileFlag == nil {
		t.Errorf("rootCmd missing 'config' persistent flag")

		return // Exit early to avoid dereferencing nil
	}

	if configFileFlag.Value.String() != config.ConfigFile {
		t.Errorf(
			"rootCmd 'config' flag value = %q, want %q",
			configFileFlag.Value.String(),
			config.ConfigFile,
		)
	}
}

// Mock os.Exit.
var osExit = os.Exit
