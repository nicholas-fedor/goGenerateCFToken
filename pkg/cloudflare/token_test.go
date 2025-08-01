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
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v5/user"
	"github.com/cloudflare/cloudflare-go/v5/zones"
	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/gogeneratecftoken/pkg/cloudflare/mocks"
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
			wantToken:   "abcdefghijklmnopqrstuvwxyz1234567890ABCD",
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone-id-123", Name: "example.com"}}}, nil).
					Once()
				m.On("CreateAPIToken", mock.Anything, mock.AnythingOfType("user.TokenNewParams")).
					Return(&user.TokenNewResponse{Value: "abcdefghijklmnopqrstuvwxyz1234567890ABCD"}, nil).
					Once()
			},
		},
		{
			name:        "ListZonesError",
			serviceName: "test-service",
			zone:        "example.com",
			wantErr:     true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPIInterface(t)
			tt.setupMock(mockAPI)

			client := &Client{Client: &cloudflare.Client{}}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() { os.Stdout = oldStdout }()

			gotToken, err := GenerateToken(
				t.Context(),
				tt.serviceName,
				tt.zone,
				client,
				mockAPI,
			)

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
