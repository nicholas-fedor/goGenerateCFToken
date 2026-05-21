/*
Copyright © 2026 Nicholas Fedor <nick@nickfedor.com>

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
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"

	mock "github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/gogeneratecftoken/pkg/config/mocks"
)

// mockViper is a simple manual mock used where the testify mock adds no value
// (e.g. error-path tests that only verify os-exit / stderr behaviour).
type mockViper struct{}

func (*mockViper) SetConfigFile(_ string)                {}
func (*mockViper) AddConfigPath(_ string)                {}
func (*mockViper) SetConfigName(_ string)                {}
func (*mockViper) SetConfigType(_ string)                {}
func (*mockViper) SetEnvPrefix(_ string)                 {}
func (*mockViper) SetEnvKeyReplacer(_ *strings.Replacer) {}
func (*mockViper) AutomaticEnv()                         {}
func (*mockViper) SetDefault(_ string, _ any)            {}
func (*mockViper) ReadInConfig() error                   { return nil }
func (*mockViper) ConfigFileUsed() string                { return "" }

func TestInitConfig(t *testing.T) {
	originalInitConfigFunc := InitConfigFunc

	defer func() { InitConfigFunc = originalInitConfigFunc }()

	called := false
	InitConfigFunc = func(_ Viper) { called = true }

	InitConfig()

	if !called {
		t.Error("InitConfig did not call InitConfigFunc")
	}
}

func TestInitConfig_NilFunc(t *testing.T) {
	originalInitConfigFunc := InitConfigFunc

	defer func() { InitConfigFunc = originalInitConfigFunc }()

	viper.Reset()

	InitConfigFunc = nil

	var buf bytes.Buffer

	originalStderr := os.Stderr

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stderr = w

	defer func() { os.Stderr = originalStderr }()

	originalUserHomeDir := osUserHomeDir
	osUserHomeDir = func() (string, error) { return "", nil }

	defer func() { osUserHomeDir = originalUserHomeDir }()

	InitConfig()

	_ = w.Close()

	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy stderr: %v", err)
	}

	output := buf.String()
	v := viper.GetViper()

	if InitConfigFunc == nil {
		t.Error("Expected InitConfigFunc to be set, but it is nil")
	}

	if v.ConfigFileUsed() == "" {
		t.Log("Config name not verifiable directly; consider mock adjustment")
	}

	if output != "" {
		t.Errorf("Expected no stderr output, got '%q'", output)
	}
}

func Test_initConfig(t *testing.T) {
	m := mocks.NewMockViper(t)

	cfgFileUsed := filepath.Join("test", "config.yaml")

	m.EXPECT().ConfigFileUsed().Return(cfgFileUsed)
	m.EXPECT().ReadInConfig().Return(nil)

	homeDir := filepath.Join("home", "test")
	expectedConfigPath := filepath.Join(homeDir, AppDirName)

	m.EXPECT().AddConfigPath(expectedConfigPath)
	m.EXPECT().AddConfigPath(".")
	m.EXPECT().SetConfigName("config")
	m.EXPECT().SetConfigType("yaml")
	m.EXPECT().SetEnvPrefix("CF")
	m.EXPECT().SetEnvKeyReplacer(mock.AnythingOfType("*strings.Replacer"))
	m.EXPECT().AutomaticEnv()
	m.EXPECT().SetDefault("api_token", "")
	m.EXPECT().SetDefault("zone", "")

	originalUserHomeDir := osUserHomeDir

	defer func() { osUserHomeDir = originalUserHomeDir }()

	osUserHomeDir = func() (string, error) { return homeDir, nil }

	var buf bytes.Buffer

	originalStderr := os.Stderr

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stderr = w

	defer func() { os.Stderr = originalStderr }()

	initConfig(m)

	_ = w.Close()

	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy stderr: %v", err)
	}

	output := buf.String()

	if output != "Using config file: "+cfgFileUsed+"\n" {
		t.Errorf("Expected output 'Using config file: %s\\n', got '%q'", cfgFileUsed, output)
	}
}

func Test_initConfig_Error(t *testing.T) {
	originalUserHomeDir := osUserHomeDir
	originalExit := osExit

	defer func() {
		osUserHomeDir = originalUserHomeDir
		osExit = originalExit
	}()

	osUserHomeDir = func() (string, error) { return "", errors.New("no home dir") }

	var exitCode int

	osExit = func(code int) { exitCode = code }

	var buf bytes.Buffer

	originalStderr := os.Stderr

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stderr = w

	defer func() { os.Stderr = originalStderr }()

	initConfig(&mockViper{})

	_ = w.Close()

	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy stderr: %v", err)
	}

	output := buf.String()

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	if !strings.Contains(output, "Failed to set config file") {
		t.Errorf("Expected stderr to contain 'Failed to set config file', got '%q'", output)
	}

	if !strings.Contains(output, "no home dir") {
		t.Errorf("Expected stderr to contain 'no home dir', got '%q'", output)
	}
}

func Test_setConfigFile(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		homeDir     string
		homeDirErr  error
		expectFile  string
		expectPaths []string
		expectName  string
		expectType  string
		expectErr   bool
	}{
		{
			name:       "CustomConfigFile",
			configFile: filepath.Join("custom", "config.yaml"),
			expectFile: filepath.Join("custom", "config.yaml"),
		},
		{
			name:        "DefaultConfig",
			homeDir:     filepath.Join("home", "test"),
			expectPaths: []string{filepath.Join(filepath.Join("home", "test"), AppDirName), "."},
			expectName:  "config",
			expectType:  "yaml",
		},
		{
			name:       "HomeDirError",
			homeDirErr: errors.New("no home dir"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockViper(t)
			originalUserHomeDir := osUserHomeDir

			defer func() { osUserHomeDir = originalUserHomeDir }()

			osUserHomeDir = func() (string, error) { return tt.homeDir, tt.homeDirErr }

			ConfigFile = tt.configFile

			defer func() { ConfigFile = "" }()

			if tt.expectErr {
				err := setConfigFile(m)
				if err == nil {
					t.Error("Expected error, got none")
				} else if !strings.Contains(err.Error(), "no home dir") {
					t.Errorf("Expected error containing 'no home dir', got '%v'", err)
				}
			} else {
				if tt.expectFile != "" {
					m.EXPECT().SetConfigFile(tt.expectFile)
				}

				if len(tt.expectPaths) > 0 {
					for _, p := range tt.expectPaths {
						m.EXPECT().AddConfigPath(p)
					}
				}

				if tt.expectName != "" {
					m.EXPECT().SetConfigName(tt.expectName)
				}

				if tt.expectType != "" {
					m.EXPECT().SetConfigType(tt.expectType)
				}

				err := setConfigFile(m)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func Test_setEnv(t *testing.T) {
	m := mocks.NewMockViper(t)

	m.EXPECT().SetEnvPrefix("CF")
	m.EXPECT().SetEnvKeyReplacer(mock.AnythingOfType("*strings.Replacer"))
	m.EXPECT().AutomaticEnv()
	m.EXPECT().SetDefault("api_token", "")
	m.EXPECT().SetDefault("zone", "")

	setEnv(m)
}

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name            string
		readConfigError error
		configFileUsed  string
		expectOutput    string
	}{
		{
			name:           "ConfigFound",
			configFileUsed: filepath.Join("test", "config.yaml"),
			expectOutput:   "Using config file: " + filepath.Join("test", "config.yaml") + "\n",
		},
		{
			name:            "ConfigNotFound",
			readConfigError: viper.ConfigFileNotFoundError{},
			expectOutput:    "",
		},
		{
			name:            "ConfigReadError",
			readConfigError: errors.New("read error"),
			expectOutput:    "Error reading config file: read error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockViper(t)

			m.EXPECT().ReadInConfig().Return(tt.readConfigError).Once()

			if tt.readConfigError == nil {
				m.EXPECT().ConfigFileUsed().Return(tt.configFileUsed).Once()
			}

			var buf bytes.Buffer

			originalStderr := os.Stderr

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}

			os.Stderr = w

			defer func() { os.Stderr = originalStderr }()

			loadConfig(m)

			_ = w.Close()

			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("Failed to copy stderr: %v", err)
			}

			output := buf.String()
			if output != tt.expectOutput {
				t.Errorf("Expected output '%q', got '%q'", tt.expectOutput, output)
			}
		})
	}
}
