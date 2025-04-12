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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v4/user"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/gogeneratecftoken/pkg/cloudflare/mocks"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		apiToken string
		wantErr  bool
	}{
		{
			name:     "Success",
			apiToken: "valid-token",
		},
		{
			name:     "MissingToken",
			apiToken: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.apiToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListZones(t *testing.T) {
	tests := []struct {
		name      string
		client    APIInterface
		wantErr   bool
		setupMock func(m *mocks.MockAPIInterface)
	}{
		{
			name:   "Success",
			client: &Client{Client: &cloudflare.Client{}},
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone-id-123"}}}, nil).
					Once()
			},
		},
		{
			name:    "NilClient",
			client:  &Client{},
			wantErr: true,
		},
		{
			name:    "ListError",
			client:  &Client{Client: &cloudflare.Client{}},
			wantErr: true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(nil, errors.New("list error")).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				mockAPI := mocks.NewMockAPIInterface(t)
				tt.setupMock(mockAPI)
				tt.client = mockAPI
			}

			_, err := tt.client.ListZones(t.Context(), zones.ZoneListParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ListZones() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateAPIToken(t *testing.T) {
	tests := []struct {
		name      string
		client    APIInterface
		wantErr   bool
		setupMock func(m *mocks.MockAPIInterface)
	}{
		{
			name:   "Success",
			client: &Client{Client: &cloudflare.Client{}},
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("CreateAPIToken", mock.Anything, mock.AnythingOfType("user.TokenNewParams")).
					Return(&user.TokenNewResponse{ID: "new-token"}, nil).
					Once()
			},
		},
		{
			name:    "NilClient",
			client:  &Client{},
			wantErr: true,
		},
		{
			name:    "CreateError",
			client:  &Client{Client: &cloudflare.Client{}},
			wantErr: true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("CreateAPIToken", mock.Anything, mock.AnythingOfType("user.TokenNewParams")).
					Return(&user.TokenNewResponse{}, errors.New("create error")).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				mockAPI := mocks.NewMockAPIInterface(t)
				tt.setupMock(mockAPI)
				tt.client = mockAPI
			}

			_, err := tt.client.CreateAPIToken(t.Context(), user.TokenNewParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAPIToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListZones_SDK(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
		wantZoneID string
	}{
		{
			name:       "Success",
			response:   `{"result":[{"id":"zone-id-123","name":"example.com"}],"success":true}`,
			statusCode: http.StatusOK,
			wantZoneID: "zone-id-123",
		},
		{
			name:       "Error",
			response:   `{"success":false,"errors":[{"message":"API error"}]}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))
				}),
			)
			defer server.Close()

			client := cloudflare.NewClient(
				option.WithHTTPClient(server.Client()),
				option.WithBaseURL(server.URL),
				option.WithAPIToken("valid-token"),
			)
			wrappedClient := &Client{Client: client}

			zones, err := wrappedClient.ListZones(t.Context(), zones.ZoneListParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ListZones() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if len(zones.Result) == 0 {
					t.Fatal("ListZones() returned no zones")
				}

				if zones.Result[0].ID != tt.wantZoneID {
					t.Errorf("ListZones() zone ID = %q, want %q", zones.Result[0].ID, tt.wantZoneID)
				}
			}
		})
	}
}

func TestClient_CreateAPIToken_SDK(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		wantErr     bool
		wantTokenID string
	}{
		{
			name:        "Success",
			response:    `{"result":{"id":"new-token"},"success":true}`,
			statusCode:  http.StatusOK,
			wantTokenID: "new-token",
		},
		{
			name:       "Error",
			response:   `{"success":false,"errors":[{"message":"API error"}]}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))
				}),
			)
			defer server.Close()

			client := cloudflare.NewClient(
				option.WithHTTPClient(server.Client()),
				option.WithBaseURL(server.URL),
				option.WithAPIToken("valid-token"),
			)
			wrappedClient := &Client{Client: client}

			token, err := wrappedClient.CreateAPIToken(t.Context(), user.TokenNewParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAPIToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && token.ID != tt.wantTokenID {
				t.Errorf("CreateAPIToken() token ID = %q, want %q", token.ID, tt.wantTokenID)
			}
		})
	}
}
