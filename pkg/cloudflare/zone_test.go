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
	"testing"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/goGenerateCFToken/pkg/cloudflare/mocks"
)

func TestGetZoneID(t *testing.T) {
	tests := []struct {
		name      string
		zone      string
		wantID    string
		wantErr   bool
		setupMock func(m *mocks.MockAPIInterface)
	}{
		{
			name:   "Success",
			zone:   "example.com",
			wantID: "zone-id-123",
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone-id-123"}}}, nil).
					Once()
			},
		},
		{
			name:    "NoZones",
			zone:    "example.com",
			wantErr: true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{}}, nil).
					Once()
			},
		},
		{
			name:    "MultipleZones",
			zone:    "example.com",
			wantErr: true,
			setupMock: func(m *mocks.MockAPIInterface) {
				m.On("ListZones", mock.Anything, mock.AnythingOfType("zones.ZoneListParams")).
					Return(&pagination.V4PagePaginationArray[zones.Zone]{Result: []zones.Zone{{ID: "zone1"}, {ID: "zone2"}}}, nil).
					Once()
			},
		},
		{
			name:    "ListError",
			zone:    "example.com",
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
			client := &Client{
				Client: &cloudflare.Client{},
			}

			mockAPI := mocks.NewMockAPIInterface(t)
			tt.setupMock(mockAPI)

			gotID, err := client.GetZoneID(t.Context(), tt.zone, mockAPI)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetZoneID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && gotID != tt.wantID {
				t.Errorf("GetZoneID() = %q, want %q", gotID, tt.wantID)
			}
		})
	}
}
