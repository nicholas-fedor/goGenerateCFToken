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
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

type mockViper struct {
	configFile      string
	configName      string
	configType      string
	configPaths     []string
	envPrefix       string
	keyReplacer     *strings.Replacer
	defaults        map[string]any
	readConfigError error
	configFileUsed  string
}

func (m *mockViper) SetConfigFile(file string) {
	m.configFile = file
}

func (m *mockViper) AddConfigPath(path string) {
	m.configPaths = append(m.configPaths, path)
}

func (m *mockViper) SetConfigName(name string) {
	m.configName = name
}

func (m *mockViper) SetConfigType(typ string) {
	m.configType = typ
}

func (m *mockViper) SetEnvPrefix(prefix string) {
	m.envPrefix = prefix
}

func (m *mockViper) SetEnvKeyReplacer(r *strings.Replacer) {
	m.keyReplacer = r
}

func (m *mockViper) AutomaticEnv() {}

func (m *mockViper) SetDefault(key string, value any) {
	if m.defaults == nil {
		m.defaults = make(map[string]any)
	}

	m.defaults[key] = value
}

func (m *mockViper) ReadInConfig() error {
	return m.readConfigError
}

func (m *mockViper) ConfigFileUsed() string {
	return m.configFileUsed
}

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

	w.Close()

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
	mock := &mockViper{
		readConfigError: viper.ConfigFileNotFoundError{},
	}
	originalUserHomeDir := osUserHomeDir

	defer func() { osUserHomeDir = originalUserHomeDir }()

	osUserHomeDir = func() (string, error) { return "/home/test", nil }

	var buf bytes.Buffer

	originalStderr := os.Stderr

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stderr = w

	defer func() { os.Stderr = originalStderr }()

	initConfig(mock)

	w.Close()

	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy stderr: %v", err)
	}

	output := buf.String()

	if mock.configName != "config" {
		t.Errorf("Expected config name 'config', got '%s'", mock.configName)
	}

	if mock.configType != "yaml" {
		t.Errorf("Expected config type 'yaml', got '%s'", mock.configType)
	}

	if len(mock.configPaths) != 2 {
		t.Errorf("Expected 2 config paths, got %d", len(mock.configPaths))
	}

	expectedPath0 := "/home/test/.goGenerateCFToken"
	actualPath0 := strings.ReplaceAll(mock.configPaths[0], string(os.PathSeparator), "/")

	if actualPath0 != expectedPath0 {
		t.Errorf("Expected first path '%s', got '%s'", expectedPath0, actualPath0)
	}

	if mock.configPaths[1] != "." {
		t.Errorf("Expected second path '.', got '%s'", mock.configPaths[1])
	}

	if mock.envPrefix != "CF" {
		t.Errorf("Expected env prefix 'CF', got '%s'", mock.envPrefix)
	}

	if mock.defaults["api_token"] != "" || mock.defaults["zone"] != "" {
		t.Errorf("Unexpected defaults: %v", mock.defaults)
	}

	if output != "" {
		t.Errorf("Expected no stderr output, got '%q'", output)
	}
}

func Test_initConfig_Error(t *testing.T) {
	mock := &mockViper{}
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

	initConfig(mock)

	w.Close()

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
			configFile: "/custom/config.yaml",
			expectFile: "/custom/config.yaml",
		},
		{
			name:        "DefaultConfig",
			homeDir:     "/home/test",
			expectPaths: []string{"/home/test/.goGenerateCFToken", "."},
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
			mock := &mockViper{}
			originalUserHomeDir := osUserHomeDir

			defer func() { osUserHomeDir = originalUserHomeDir }()

			osUserHomeDir = func() (string, error) { return tt.homeDir, tt.homeDirErr }

			ConfigFile = tt.configFile
			defer func() { ConfigFile = "" }()

			err := setConfigFile(mock)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got none")
				} else if !strings.Contains(err.Error(), "no home dir") {
					t.Errorf("Expected error containing 'no home dir', got '%v'", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if tt.expectFile != "" && mock.configFile != tt.expectFile {
					t.Errorf("Expected config file '%s', got '%s'", tt.expectFile, mock.configFile)
				}

				if len(tt.expectPaths) > 0 && len(mock.configPaths) != len(tt.expectPaths) {
					t.Errorf("Expected %d paths, got %d", len(tt.expectPaths), len(mock.configPaths))
				}

				for i, path := range tt.expectPaths {
					if i < len(mock.configPaths) {
						actualPath := strings.ReplaceAll(mock.configPaths[i], string(os.PathSeparator), "/")
						if actualPath != path {
							t.Errorf("Expected path %d to be '%s', got '%s'", i, path, actualPath)
						}
					}
				}

				if tt.expectName != "" && mock.configName != tt.expectName {
					t.Errorf("Expected config name '%s', got '%s'", tt.expectName, mock.configName)
				}

				if tt.expectType != "" && mock.configType != tt.expectType {
					t.Errorf("Expected config type '%s', got '%s'", tt.expectType, mock.configType)
				}
			}
		})
	}
}

func Test_setEnv(t *testing.T) {
	mock := &mockViper{}
	setEnv(mock)

	if mock.envPrefix != "CF" {
		t.Errorf("Expected env prefix 'CF', got '%s'", mock.envPrefix)
	}

	if mock.keyReplacer == nil || mock.keyReplacer.Replace("a.b") != "a_b" {
		t.Errorf("Expected key replacer to replace '.' with '_', got %v", mock.keyReplacer)
	}

	if mock.defaults["api_token"] != "" || mock.defaults["zone"] != "" {
		t.Errorf("Unexpected defaults: %v", mock.defaults)
	}
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
			configFileUsed: "/test/config.yaml",
			expectOutput:   "Using config file: /test/config.yaml\n",
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
			mock := &mockViper{
				readConfigError: tt.readConfigError,
				configFileUsed:  tt.configFileUsed,
			}

			var buf bytes.Buffer

			originalStderr := os.Stderr

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}

			os.Stderr = w

			defer func() { os.Stderr = originalStderr }()

			loadConfig(mock)

			w.Close()

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
