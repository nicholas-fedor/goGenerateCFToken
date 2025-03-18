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
package cloudflare

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"
)

// mockAPI implements APIInterface for testing.
type mockAPI struct {
	listZones   func(ctx context.Context, zone ...string) ([]cloudflare.Zone, error)
	createToken func(ctx context.Context, token cloudflare.APIToken) (cloudflare.APIToken, error)
}

func (m *mockAPI) ListZones(ctx context.Context, zone ...string) ([]cloudflare.Zone, error) {
	return m.listZones(ctx, zone...)
}

func (m *mockAPI) CreateAPIToken(ctx context.Context, token cloudflare.APIToken) (cloudflare.APIToken, error) {
	return m.createToken(ctx, token)
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		zone        string
		listZones   func(ctx context.Context, zone ...string) ([]cloudflare.Zone, error)
		createToken func(ctx context.Context, token cloudflare.APIToken) (cloudflare.APIToken, error)
		wantToken   string
		wantErr     bool
	}{
		{
			name:        "Success",
			serviceName: "test-service",
			zone:        "example.com",
			listZones: func(_ context.Context, zone ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{{ID: "zone-id-123", Name: "example.com"}}, nil
			},
			createToken: func(_ context.Context, _ cloudflare.APIToken) (cloudflare.APIToken, error) {
				return cloudflare.APIToken{Value: "new-token"}, nil
			},
			wantToken: "new-token",
		},
		{
			name:        "ListZonesError",
			serviceName: "test-service",
			zone:        "example.com",
			listZones: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return nil, errors.New("zone error")
			},
			createToken: func(_ context.Context, _ cloudflare.APIToken) (cloudflare.APIToken, error) {
				return cloudflare.APIToken{}, nil
			},
			wantErr: true,
		},
		{
			name:        "CreateTokenError",
			serviceName: "test-service",
			zone:        "example.com",
			listZones: func(_ context.Context, zone ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{{ID: "zone-id-123", Name: "example.com"}}, nil
			},
			createToken: func(_ context.Context, _ cloudflare.APIToken) (cloudflare.APIToken, error) {
				return cloudflare.APIToken{}, errors.New("create error")
			},
			wantErr: true,
		},
		{
			name:        "NoZones",
			serviceName: "test-service",
			zone:        "example.com",
			listZones: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{}, nil
			},
			createToken: func(_ context.Context, _ cloudflare.APIToken) (cloudflare.APIToken, error) {
				return cloudflare.APIToken{}, nil
			},
			wantErr: true,
		},
		{
			name:        "MultipleZones",
			serviceName: "test-service",
			zone:        "example.com",
			listZones: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{{ID: "zone1"}, {ID: "zone2"}}, nil
			},
			createToken: func(_ context.Context, _ cloudflare.APIToken) (cloudflare.APIToken, error) {
				return cloudflare.APIToken{}, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &mockAPI{
				listZones:   tt.listZones,
				createToken: tt.createToken,
			}

			// Capture stdout for token name output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() { os.Stdout = oldStdout }()

			gotToken, err := GenerateToken(tt.serviceName, tt.zone, api, context.Background())

			w.Close()

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if gotToken != tt.wantToken {
					t.Errorf("GenerateToken() = %q, want %q", gotToken, tt.wantToken)
				}

				if output != "Generating API token: "+tt.serviceName+"."+tt.zone+"\n" {
					t.Errorf("GenerateToken() output = %q, want %q", output, "Generating API token: "+tt.serviceName+"."+tt.zone+"\n")
				}
			}
		})
	}
}
