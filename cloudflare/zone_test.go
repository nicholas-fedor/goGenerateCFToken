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
	"testing"

	"github.com/cloudflare/cloudflare-go"
)

// mockZoneAPIForGetZoneID implements APIInterface for testing GetZoneID.
type mockZoneAPIForGetZoneID struct {
	listZones func(ctx context.Context, zone ...string) ([]cloudflare.Zone, error)
}

func (m *mockZoneAPIForGetZoneID) ListZones(ctx context.Context, zone ...string) ([]cloudflare.Zone, error) {
	return m.listZones(ctx, zone...)
}

// CreateAPIToken is required by APIInterface but not used in GetZoneID tests.
func (m *mockZoneAPIForGetZoneID) CreateAPIToken(ctx context.Context, token cloudflare.APIToken) (cloudflare.APIToken, error) {
	return cloudflare.APIToken{}, nil // No-op for GetZoneID tests
}

func TestGetZoneID(t *testing.T) {
	tests := []struct {
		name     string
		zone     string
		listFunc func(ctx context.Context, zone ...string) ([]cloudflare.Zone, error)
		wantID   string
		wantErr  bool
	}{
		{
			name: "Success",
			zone: "example.com",
			listFunc: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{{ID: "zone-id-123"}}, nil
			},
			wantID: "zone-id-123",
		},
		{
			name: "NoZones",
			zone: "example.com",
			listFunc: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{}, nil
			},
			wantErr: true,
		},
		{
			name: "MultipleZones",
			zone: "example.com",
			listFunc: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return []cloudflare.Zone{{ID: "zone1"}, {ID: "zone2"}}, nil
			},
			wantErr: true,
		},
		{
			name: "ListError",
			zone: "example.com",
			listFunc: func(_ context.Context, _ ...string) ([]cloudflare.Zone, error) {
				return nil, errors.New("list error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &mockZoneAPIForGetZoneID{
				listZones: tt.listFunc,
			}

			ctx := context.Background()

			gotID, err := GetZoneID(tt.zone, api, ctx)
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
