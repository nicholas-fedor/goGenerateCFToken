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

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/config"
)

func TestRootCmd(t *testing.T) {
	if rootCmd.Use != "goGenerateCFToken" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "goGenerateCFToken")
	}

	if rootCmd.Version != version {
		t.Errorf("rootCmd.Version = %q, want %q", rootCmd.Version, version)
	}

	if rootCmd.Short == "" || rootCmd.Long == "" {
		t.Errorf("rootCmd Short or Long description is empty")
	}

	oldExit := osExit
	osExit = func(code int) {
		if code != 0 {
			t.Errorf("Execute() exited with code %d", code)
		}
	}

	defer func() { osExit = oldExit }()

	oldConfigFile := config.ConfigFile
	config.ConfigFile = "dummy.yaml"

	defer func() { config.ConfigFile = oldConfigFile }()

	Execute()

	configFileFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFileFlag == nil {
		t.Errorf("rootCmd missing 'config' persistent flag")

		return
	}

	if configFileFlag.Value.String() != config.ConfigFile {
		t.Errorf(
			"rootCmd 'config' flag value = %q, want %q",
			configFileFlag.Value.String(),
			config.ConfigFile,
		)
	}
}

var osExit = os.Exit
