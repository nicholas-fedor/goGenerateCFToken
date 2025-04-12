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
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/cloudflare/cloudflare-go/v4/user"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/cloudflare/mocks"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		zone        string
		zoneID      string
		wantToken   string
		wantErr     bool
		setupMock   func(m *mocks.MockAPIInterface)
	}{
		{
			name:        "Success",
			serviceName: "test-service",
			zone:        "example.com",
			zoneID:      "zone-id-123",
			wantToken:   "new-token",
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone-id-123", Name: "example.com"}}}, nil).
					Once()
				m.On("CreateAPIToken", mock.Anything, mock.AnythingOfType("user.TokenNewParams")).
					Return(&user.TokenNewResponse{ID: "new-token"}, nil).
					Once()
			},
		},
		{
			name:        "ListZonesError",
			serviceName: "test-service",
			zone:        "example.com",
			wantErr:     true,
			setupMock: func(_m *mocks.MockAPIInterface) {
				_m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(nil, errors.New("zone error")).
					Once()
			},
		},
		{
			name:        "CreateTokenError",
			serviceName: "test-service",
			zone:        "example.com",
			zoneID:      "zone-id-123",
			wantErr:     true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone-id-123", Name: "example.com"}}}, nil).
					Once()
				m.On("CreateAPIToken", mock.Anything, mock.AnythingOfType("user.TokenNewParams")).
					Return(&user.TokenNewResponse{}, errors.New("create error")).
					Once()
			},
		},
		{
			name:        "NoZones",
			serviceName: "test-service",
			zone:        "example.com",
			wantErr:     true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{}}, nil).
					Once()
			},
		},
		{
			name:        "MultipleZones",
			serviceName: "test-service",
			zone:        "example.com",
			wantErr:     true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone1"}, {ID: "zone2"}}}, nil).
					Once()
			},
		},
		{
			name:        "NilClient",
			serviceName: "test-service",
			zone:        "example.com",
			wantErr:     true,
			setupMock: func(_m *mocks.MockAPIInterface) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPIInterface(t)
			tt.setupMock(mockAPI)

			var client *Client
			if tt.name != "NilClient" {
				client = &Client{
					Client: &cloudflare.Client{},
				}
			}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() { os.Stdout = oldStdout }()

			var gotToken string

			var err error
			if client != nil {
				gotToken, err = generateTokenForTest(
					t.Context(),
					tt.serviceName,
					tt.zone,
					mockAPI,
				)
			} else {
				_, err = (&Client{}).GenerateToken(t.Context(), tt.serviceName, tt.zone)
			}

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
					t.Errorf(
						"GenerateToken() output = %q, want %q",
						output,
						"Generating API token: "+tt.serviceName+"."+tt.zone+"\n",
					)
				}
			}
		})
	}
}

func generateTokenForTest(
	ctx context.Context,
	serviceName, zoneName string,
	api APIInterface,
) (string, error) {
	client := &Client{}

	zoneID, err := client.GetZoneID(ctx, zoneName, api)
	if err != nil {
		return "", fmt.Errorf("failed to get zone ID: %w", err)
	}

	tokenName := serviceName + "." + zoneName
	fmt.Fprintln(os.Stdout, "Generating API token:", tokenName)
	params := user.TokenNewParams{
		Name: cloudflare.F(tokenName),
		Policies: cloudflare.F([]shared.TokenPolicyParam{{
			Effect: cloudflare.F(shared.TokenPolicyEffectAllow),
			PermissionGroups: cloudflare.F([]shared.TokenPolicyPermissionGroupParam{{
				ID: cloudflare.F("c8fed203ed3043cba015a93ad1616f1f"),
			}, {
				ID: cloudflare.F("4755a26eedb94da69e1066d98aa820be"),
			}}),
			Resources: cloudflare.F(map[string]string{
				"com.cloudflare.api.account.zone." + zoneID: "*",
			}),
		}}),
	}

	token, err := api.CreateAPIToken(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token.ID, nil
}

func TestClientAdapter(t *testing.T) {
	tests := []struct {
		name         string
		client       *cloudflare.Client
		wantListErr  string
		wantTokenErr string
	}{
		{
			name:         "NilClient",
			client:       nil,
			wantListErr:  ErrZonesServiceNotInit.Error(),
			wantTokenErr: ErrTokensServiceNotInit.Error(),
		},
		{
			name:         "EmptyClient",
			client:       &cloudflare.Client{},
			wantListErr:  ErrZonesServiceNotInit.Error(),
			wantTokenErr: ErrTokensServiceNotInit.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &clientAdapter{tt.client}

			// Test ListZones
			_, err := adapter.ListZones(t.Context(), zones.ZoneListParams{})
			if err == nil || err.Error() != tt.wantListErr {
				t.Errorf("ListZones() error = %v, want %q", err, tt.wantListErr)
			}

			// Test CreateAPIToken
			_, err = adapter.CreateAPIToken(t.Context(), user.TokenNewParams{})
			if err == nil || err.Error() != tt.wantTokenErr {
				t.Errorf("CreateAPIToken() error = %v, want %q", err, tt.wantTokenErr)
			}
		})
	}
}
