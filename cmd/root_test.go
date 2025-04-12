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
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/config"
)

var osExit = os.Exit

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

	var exitCode int

	osExit = func(code int) {
		exitCode = code
	}

	defer func() { osExit = os.Exit }()

	oldConfigFile := config.ConfigFile
	config.ConfigFile = "dummy.yaml"

	defer func() { config.ConfigFile = oldConfigFile }()

	Execute()

	if exitCode != 0 {
		t.Errorf("Execute() exited with code %d, want 0", exitCode)
	}

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

func TestUserHomeDir(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		wantPath string
	}{
		{
			name:     "Windows",
			goos:     "windows",
			wantPath: "%userprofile%",
		},
		{
			name:     "Linux",
			goos:     "linux",
			wantPath: "$HOME",
		},
		{
			name:     "iOS",
			goos:     "ios",
			wantPath: "/",
		},
		{
			name:     "Plan9",
			goos:     "plan9",
			wantPath: "$home",
		},
		{
			name:     "Android",
			goos:     "android",
			wantPath: "/sdcard",
		},
		{
			name:     "Unknown",
			goos:     "unknown",
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origGOOS := goos
			goos = tt.goos

			defer func() { goos = origGOOS }()

			got := userHomeDir()
			if got != tt.wantPath {
				t.Errorf("userHomeDir() = %q, want %q", got, tt.wantPath)
			}
		})
	}
}

func TestConfigPath(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		wantPath string
	}{
		{
			name:     "Windows",
			goos:     "windows",
			wantPath: filepath.Join("%userprofile%", ".goGenerateCFToken", "config.yaml"),
		},
		{
			name:     "Linux",
			goos:     "linux",
			wantPath: filepath.Join("$HOME", ".goGenerateCFToken", "config.yaml"),
		},
		{
			name:     "Unknown",
			goos:     "unknown",
			wantPath: filepath.Join("", ".goGenerateCFToken", "config.yaml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origGOOS := goos
			goos = tt.goos

			defer func() { goos = origGOOS }()

			got := configPath()
			if got != tt.wantPath {
				t.Errorf("configPath() = %q, want %q", got, tt.wantPath)
			}
		})
	}
}

func TestRootCmd_ConfigFlag(t *testing.T) {
	origConfigFile := config.ConfigFile
	defer func() { config.ConfigFile = origConfigFile }()

	config.ConfigFile = "test-config.yaml"
	rootCmd := &cobra.Command{Use: "goGenerateCFToken"}
	rootCmd.PersistentFlags().StringVar(&config.ConfigFile, "config", "", "config file")

	if err := rootCmd.ParseFlags([]string{"--config", "custom-config.yaml"}); err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	if config.ConfigFile != "custom-config.yaml" {
		t.Errorf("config.ConfigFile = %q, want %q", config.ConfigFile, "custom-config.yaml")
	}

	configFileFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFileFlag.Value.String() != "custom-config.yaml" {
		t.Errorf(
			"config flag value = %q, want %q",
			configFileFlag.Value.String(),
			"custom-config.yaml",
		)
	}
}
