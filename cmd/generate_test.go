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
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"                                  // For cloudflare.API
	cloudflarePkg "github.com/nicholas-fedor/goGenerateCFToken/cloudflare" // Alias to avoid conflict
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MockAPI implements cloudflarePkg.APIInterface.
type MockAPI struct {
	cloudflarePkg.APIInterface
}

func TestGenerateCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		apiToken   string
		zone       string
		clientFunc func(token string) (*cloudflare.API, error)
		genFunc    func(serviceName, zone string, api cloudflarePkg.APIInterface, ctx context.Context) (string, error)
		wantErr    bool
		wantOutput string
	}{
		{
			name:     "Success",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(token string) (*cloudflare.API, error) {
				return &cloudflare.API{}, nil
			},
			genFunc: func(_, _ string, _ cloudflarePkg.APIInterface, _ context.Context) (string, error) {
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
			clientFunc: func(token string) (*cloudflare.API, error) {
				return nil, fmt.Errorf("client error")
			},
			wantErr: true,
		},
		{
			name:     "GenerateError",
			args:     []string{"generate", "test-service"},
			apiToken: "valid-token",
			zone:     "example.com",
			clientFunc: func(token string) (*cloudflare.API, error) {
				return &cloudflare.API{}, nil
			},
			genFunc: func(_, _ string, _ cloudflarePkg.APIInterface, _ context.Context) (string, error) {
				return "", fmt.Errorf("generate error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test
			viper.Reset()
			viper.Set("api_token", tt.apiToken)
			viper.Set("zone", tt.zone)

			// Save and restore original functions
			origNewAPIClient := NewAPIClientFunc
			origGenerateToken := GenerateTokenFunc

			defer func() {
				NewAPIClientFunc = origNewAPIClient
				GenerateTokenFunc = origGenerateToken
			}()

			// Apply mocks explicitly
			if tt.clientFunc != nil {
				NewAPIClientFunc = tt.clientFunc
			} else {
				NewAPIClientFunc = func(token string) (*cloudflare.API, error) {
					return &cloudflare.API{}, nil
				}
			}

			if tt.genFunc != nil {
				GenerateTokenFunc = tt.genFunc
			} else {
				GenerateTokenFunc = func(_, _ string, _ cloudflarePkg.APIInterface, _ context.Context) (string, error) {
					return "new-token", nil
				}
			}

			// Create a fresh rootCmd and add the real generateCmd
			rootCmd := &cobra.Command{Use: "goGenerateCFToken"}
			// Reset flags and rebind to avoid interference
			generateCmd.ResetFlags()
			generateCmd.Flags().StringP("token", "t", "", "Cloudflare API token")
			generateCmd.Flags().StringP("zone", "z", "", "Cloudflare zone name")
			viper.BindPFlag("api_token", generateCmd.Flags().Lookup("token"))
			viper.BindPFlag("zone", generateCmd.Flags().Lookup("zone"))
			rootCmd.AddCommand(generateCmd)

			// Capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() { os.Stdout = oldStdout }()

			// Execute the command
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
