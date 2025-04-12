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
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/cloudflare"
	"github.com/nicholas-fedor/goGenerateCFToken/pkg/config"
)

func TestGenerateCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		apiToken   string
		zone       string
		clientFunc func(apiToken string) (*cloudflare.Client, error)
		genFunc    func(c *cloudflare.Client, ctx context.Context, serviceName string, zone string) (string, error)
		wantErr    bool
		wantOutput string
	}{
		{
			name:     "Success",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return &cloudflare.Client{}, nil
			},
			genFunc: func(_ *cloudflare.Client, _ context.Context, _ string, _ string) (string, error) {
				return "new-token", nil
			},
			wantOutput: "new-token\n",
		},
		{
			name:    "MissingArgs",
			args:    []string{"generate"},
			wantErr: true,
		},
		{
			name:    "MissingAPIToken",
			args:    []string{"generate", "test-service"},
			zone:    "example.com",
			wantErr: true,
		},
		{
			name:     "MissingZone",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			wantErr:  true,
		},
		{
			name:     "ClientError",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return nil, errors.New("client error")
			},
			wantErr: true,
		},
		{
			name:     "GenerateError",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(_ string) (*cloudflare.Client, error) {
				return &cloudflare.Client{}, nil
			},
			genFunc: func(_ *cloudflare.Client, _ context.Context, _ string, _ string) (string, error) {
				return "", errors.New("generate error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set("api_token", tt.apiToken)
			viper.Set("zone", tt.zone)

			origInitConfig := config.InitConfigFunc
			defer func() { config.InitConfigFunc = origInitConfig }()

			config.InitConfigFunc = func() {
				viper.Set("api_token", tt.apiToken)
				viper.Set("zone", tt.zone)
			}

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
				GenerateTokenFunc = func(_ *cloudflare.Client, _ context.Context, _ string, _ string) (string, error) {
					return "new-token", nil
				}
			}

			rootCmd := &cobra.Command{Use: "goGenerateCFToken"}

			generateCmd.ResetFlags()
			generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
			generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")
			viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token"))
			viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone"))
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
		})
	}
}
