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
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/cloudflare"
	"github.com/nicholas-fedor/goGenerateCFToken/pkg/config"
)

var newToken = "new-token"

func TestGenerateCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		apiToken   string
		zone       string
		clientFunc func(apiToken string) (*cloudflare.Client, error)
		genFunc    func(ctx context.Context, serviceName string, zone string, client *cloudflare.Client, api cloudflare.APIInterface) (string, error)
		configFile string
		configErr  bool
		wantErr    bool
		wantOutput string
		wantErrMsg string
	}{
		{
			name:     "Success",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return &cloudflare.Client{}, nil
			},
			genFunc: func(_ context.Context, _ string, _ string, _ *cloudflare.Client, _ cloudflare.APIInterface) (string, error) {
				return newToken, nil
			},
			wantOutput: "new-token\n",
		},
		{
			name:       "MissingArgs",
			args:       []string{"generate"},
			wantErr:    true,
			wantErrMsg: "accepts 1 arg(s), received 0",
		},
		{
			name:       "MissingAPIToken",
			args:       []string{"generate", "test-service"},
			zone:       "example.com",
			wantErr:    true,
			wantErrMsg: "missing required credentials in config",
		},
		{
			name:       "MissingZone",
			args:       []string{"generate", "test-service"},
			apiToken:   "valid-token",
			wantErr:    true,
			wantErrMsg: "missing required zone in config",
		},
		{
			name:     "ClientError",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return nil, errors.New("client error")
			},
			wantErr:    true,
			wantErrMsg: "failed to initialize Cloudflare client: client error",
		},
		{
			name:     "GenerateError",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return &cloudflare.Client{}, nil
			},
			genFunc: func(_ context.Context, _ string, _ string, _ *cloudflare.Client, _ cloudflare.APIInterface) (string, error) {
				return "", errors.New("generate error")
			},
			wantErr:    true,
			wantErrMsg: "failed to generate token: generate error",
		},
		{
			name:       "InvalidFlag",
			args:       []string{"generate", "test-service", "--invalid"},
			apiToken:   "valid-token",
			zone:       "example.com",
			wantErr:    true,
			wantErrMsg: "unknown flag: --invalid",
		},
		{
			name:     "MalformedServiceName",
			args:     []string{"generate", "test@service#invalid"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return &cloudflare.Client{}, nil
			},
			genFunc: func(_ context.Context, serviceName string, _ string, _ *cloudflare.Client, _ cloudflare.APIInterface) (string, error) {
				if strings.ContainsAny(serviceName, "@#") {
					return "", errors.New("invalid service name")
				}

				return newToken, nil
			},
			wantErr:    true,
			wantErrMsg: "failed to generate token: invalid service name",
		},
		{
			name:     "ConflictingFlagToken",
			args:     []string{"generate", "test-service", "--token", "flag-token"},
			apiToken: "config-token",
			zone:     "example.com",
			clientFunc: func(token string) (*cloudflare.Client, error) {
				if token != "flag-token" {
					return nil, errors.New("expected flag-token")
				}

				return &cloudflare.Client{}, nil
			},
			genFunc: func(_ context.Context, _ string, _ string, _ *cloudflare.Client, _ cloudflare.APIInterface) (string, error) {
				return newToken, nil
			},
			wantOutput: "new-token\n",
		},
		{
			name:       "InvalidConfigFile",
			args:       []string{"generate", "test-service"},
			configFile: "invalid.yaml",
			configErr:  true,
			wantErr:    true,
			wantErrMsg: "missing required credentials in config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			origInitConfig := config.InitConfigFunc
			defer func() { config.InitConfigFunc = origInitConfig }()

			config.InitConfigFunc = func(v config.Viper) {
				if tt.configErr {
					fmt.Fprintf(
						os.Stderr,
						"Error reading config file: open %s: no such file or directory\n",
						tt.configFile,
					)

					return
				}

				v.SetDefault("api_token", tt.apiToken)
				v.SetDefault("zone", tt.zone)
			}

			origConfigFile := config.ConfigFile
			defer func() { config.ConfigFile = origConfigFile }()

			config.ConfigFile = tt.configFile

			origNewClient := NewClientFunc
			origGenerateToken := GenerateTokenFunc

			defer func() {
				NewClientFunc = origNewClient
				GenerateTokenFunc = origGenerateToken
			}()

			if tt.clientFunc != nil {
				NewClientFunc = tt.clientFunc
			} else {
				NewClientFunc = func(_ string) (*cloudflare.Client, error) {
					return &cloudflare.Client{}, nil
				}
			}

			if tt.genFunc != nil {
				GenerateTokenFunc = tt.genFunc
			} else {
				GenerateTokenFunc = func(_ context.Context, _ string, _ string, _ *cloudflare.Client, _ cloudflare.APIInterface) (string, error) {
					return newToken, nil
				}
			}

			rootCmd := &cobra.Command{Use: "goGenerateCFToken"}

			generateCmd.ResetFlags()
			generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
			generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")

			if err := BindPFlagFunc("api_token", generateCmd.Flags().Lookup("token")); err != nil {
				t.Fatalf("Failed to bind api_token: %v", err)
			}

			if err := BindPFlagFunc("zone", generateCmd.Flags().Lookup("zone")); err != nil {
				t.Fatalf("Failed to bind zone: %v", err)
			}

			rootCmd.AddCommand(generateCmd)

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() { os.Stdout = oldStdout }()

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()

			w.Close()

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if (err != nil) != tt.wantErr {
				t.Errorf("rootCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && output != tt.wantOutput {
				t.Errorf("rootCmd.Execute() output = %q, want %q", output, tt.wantOutput)
			}

			if tt.wantErr && tt.wantErrMsg != "" &&
				(err == nil || !strings.Contains(err.Error(), tt.wantErrMsg)) {
				t.Errorf("rootCmd.Execute() error = %v, wantErrMsg %q", err, tt.wantErrMsg)
			}
		})
	}
}

func TestGenerateCmd_BindErrors(t *testing.T) {
	tests := []struct {
		name      string
		flag      string
		wantPanic bool
	}{
		{
			name:      "BindAPITokenFlagError",
			flag:      "api_token",
			wantPanic: true,
		},
		{
			name:      "BindZoneFlagError",
			flag:      "zone",
			wantPanic: true,
		},
		{
			name:      "BindAPITokenFlagErrorMock",
			flag:      "api_token",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			origBindPFlag := BindPFlagFunc
			defer func() { BindPFlagFunc = origBindPFlag }()

			if tt.name == "BindAPITokenFlagErrorMock" {
				BindPFlagFunc = func(key string, flag *pflag.Flag) error {
					if key == tt.flag {
						return errors.New("bind error")
					}

					return viper.BindPFlag(key, flag)
				}
			}

			rootCmd := &cobra.Command{Use: "goGenerateCFToken"}

			generateCmd.ResetFlags()

			if tt.flag != "api_token" {
				generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
			}

			if tt.flag != "zone" {
				generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")
			}

			rootCmd.AddCommand(generateCmd)

			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("Expected panic = %v, got %v", tt.wantPanic, r)
				}
			}()

			if err := viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token")); err != nil {
				panic(fmt.Errorf("%w: %w", ErrBindAPITokenFlag, err))
			}

			if err := viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone")); err != nil {
				panic(fmt.Errorf("%w: %w", ErrBindZoneFlag, err))
			}
		})
	}
}
