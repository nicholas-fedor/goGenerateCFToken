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
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	// Reset viper to ensure a clean state for each test
	viper.Reset()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file
	configContent := []byte("api_token: test-api-token\nzone: example.com")
	configPath := filepath.Join(tempDir, "config.yaml")

	if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	tests := []struct {
		name       string
		configFile string
		envVars    map[string]string
		wantToken  string
		wantZone   string
		wantErr    bool
	}{
		{
			name:       "ConfigFileSpecified",
			configFile: configPath,
			wantToken:  "test-api-token",
			wantZone:   "example.com",
		},
		{
			name:       "NoConfigFileWithEnvVars",
			configFile: "",
			envVars: map[string]string{
				"CF_API_TOKEN": "env-api-token",
				"CF_ZONE":      "env-example.com",
			},
			wantToken: "env-api-token",
			wantZone:  "env-example.com",
		},
		{
			name:       "DefaultConfigNotFound",
			configFile: "",
			wantToken:  "",
			wantZone:   "",
		},
		{
			name:       "InvalidConfigFile",
			configFile: filepath.Join(tempDir, "nonexistent.yaml"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test case
			viper.Reset()

			// Set environment variables if provided
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Set ConfigFile for this test
			ConfigFile = tt.configFile

			// Redirect stderr to capture output
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			defer func() { os.Stderr = oldStderr }()

			InitConfig()

			// Check viper values
			if gotToken := viper.GetString("api_token"); gotToken != tt.wantToken {
				t.Errorf("InitConfig() api_token = %v, want %v", gotToken, tt.wantToken)
			}

			if gotZone := viper.GetString("zone"); gotZone != tt.wantZone {
				t.Errorf("InitConfig() zone = %v, want %v", gotZone, tt.wantZone)
			}

			// Check stderr output for errors
			w.Close()

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.wantErr && !containsError(output) {
				t.Errorf("InitConfig() expected error in output, got: %s", output)
			}
		})
	}
}

// Helper function to check if output contains an error message.
func containsError(output string) bool {
	return len(output) > 0 && output != "Using config file: "+ConfigFile+"\n"
}
