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
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		envVars    map[string]string
		wantToken  string
		wantZone   string
		wantOutput string
		wantErrOut bool
		useMock    bool
	}{
		{
			name:       "ConfigFileSpecified",
			configFile: "config.yaml",
			wantToken:  "test-api-token",
			wantZone:   "example.com",
			wantOutput: "",
			useMock:    false,
		},
		{
			name: "NoConfigFileWithEnvVars",
			envVars: map[string]string{
				"CF_API_TOKEN": "env-api-token",
				"CF_ZONE":      "env-example.com",
			},
			wantToken:  "env-api-token",
			wantZone:   "env-example.com",
			wantOutput: "",
			useMock:    true,
		},
		{
			name:       "DefaultConfigNotFound",
			configFile: "",
			wantToken:  "",
			wantZone:   "",
			wantOutput: "",
			useMock:    true,
		},
		{
			name:       "InvalidConfigFile",
			configFile: "invalid.yaml",
			wantErrOut: true,
			useMock:    false,
		},
		{
			name:       "HomeDirError",
			configFile: "",
			wantErrOut: true,
			useMock:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			if tt.configFile != "" && !tt.wantErrOut {
				configContent := []byte("api_token: test-api-token\nzone: example.com")
				configPath := filepath.Join(tempDir, "config.yaml")

				if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
					t.Fatalf("Failed to write config file: %v", err)
				}

				tt.configFile = configPath
				tt.wantOutput = "Using config file: " + configPath + "\n"
			} else if tt.configFile != "" {
				tt.configFile = filepath.Join(tempDir, tt.configFile)
			}

			viper.Reset()

			os.Unsetenv("CF_API_TOKEN")
			os.Unsetenv("CF_ZONE")

			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			ConfigFile = tt.configFile

			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			defer func() { os.Stderr = oldStderr }()

			if tt.useMock {
				origInitConfig := InitConfigFunc
				defer func() { InitConfigFunc = origInitConfig }()

				InitConfigFunc = func() {
					if tt.name == "HomeDirError" {
						fmt.Fprintf(os.Stderr, "Error: no home directory\n")

						return
					}

					if ConfigFile != "" {
						viper.SetConfigFile(ConfigFile)
					} else {
						viper.AddConfigPath(tempDir)
						viper.SetConfigName("nonexistent")
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
			}

			InitConfig()

			w.Close()

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if gotToken := viper.GetString("api_token"); gotToken != tt.wantToken {
				t.Errorf("InitConfig() api_token = %q, want %q", gotToken, tt.wantToken)
			}

			if gotZone := viper.GetString("zone"); gotZone != tt.wantZone {
				t.Errorf("InitConfig() zone = %q, want %q", gotZone, tt.wantZone)
			}

			if tt.wantErrOut && output == "" {
				t.Errorf("InitConfig() expected error output, got none")
			}

			if !tt.wantErrOut && output != tt.wantOutput {
				t.Errorf("InitConfig() output = %q, want %q", output, tt.wantOutput)
			}
		})
	}
}
