// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/shared"
	"github.com/cloudflare/cloudflare-go/v7/user"
	"github.com/cloudflare/cloudflare-go/v7/zones"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/cloudflare/mocks"
)

func Test_resolveZoneID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		zoneName string
		mockFn   func(m *mocks.MockAPI)
		want     string
		wantErr  bool
		errType  error
	}{
		{
			name:     "zone found",
			zoneName: "example.com",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{
						{ID: "zone123", Name: "example.com"},
					}, nil)
			},
			want:    "zone123",
			wantErr: false,
		},
		{
			name:     "zone not found",
			zoneName: "missing.com",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{}, nil)
			},
			want:    "",
			wantErr: true,
			errType: ErrZoneNotFound,
		},
		{
			name:     "multiple zones found",
			zoneName: "duplicate.com",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{
						{ID: "zone1", Name: "duplicate.com"},
						{ID: "zone2", Name: "duplicate.com"},
					}, nil)
			},
			want:    "",
			wantErr: true,
			errType: ErrMultipleZonesFound,
		},
		{
			name:     "api error propagates",
			zoneName: "example.com",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return(nil, errors.New("api failure"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPI(t)
			tt.mockFn(mockAPI)

			got, err := resolveZoneID(ctx, mockAPI, tt.zoneName, "")

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, got)

				if tt.errType != nil {
					require.ErrorIs(t, err, tt.errType)
				}

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_buildTokenPolicy(t *testing.T) {
	tests := []struct {
		name   string
		zoneID string
	}{
		{
			name:   "builds policy for zone",
			zoneID: "abc123",
		},
		{
			name:   "different zone ID",
			zoneID: "xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildTokenPolicy(tt.zoneID)
			require.Len(t, got, 1)
			assert.Equal(t, cloudflare.F(shared.TokenPolicyEffectAllow), got[0].Effect)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		req     TokenGenerationRequest
		mockFn  func(m *mocks.MockAPI)
		want    *TokenGenerationResult
		wantErr bool
		errType error
	}{
		{
			name: "successful token generation",
			req: TokenGenerationRequest{
				ServiceName: "myservice",
				ZoneName:    "example.com",
			},
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{{ID: "zone123", Name: "example.com"}}, nil)
				m.EXPECT().CreateAPIToken(mock.Anything, mock.Anything).
					Return(&user.TokenNewResponse{
						Value: "generated-token-value",
						ID:    "new-token-id",
					}, nil)
			},
			want: &TokenGenerationResult{
				Token: "generated-token-value",
				Name:  "myservice.example.com",
				Zone:  "example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid service name",
			req: TokenGenerationRequest{
				ServiceName: "invalid service!",
				ZoneName:    "example.com",
			},
			mockFn:  func(m *mocks.MockAPI) {},
			wantErr: true,
			errType: ErrInvalidServiceName,
		},
		{
			name: "empty service name",
			req: TokenGenerationRequest{
				ServiceName: "",
				ZoneName:    "example.com",
			},
			mockFn:  func(m *mocks.MockAPI) {},
			wantErr: true,
			errType: ErrInvalidServiceName,
		},
		{
			name: "custom name overrides default",
			req: TokenGenerationRequest{
				ServiceName: "myservice",
				ZoneName:    "example.com",
				Name:        "custom-token-name",
			},
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{{ID: "zone123", Name: "example.com"}}, nil)
				m.EXPECT().CreateAPIToken(mock.Anything, mock.Anything).
					Return(&user.TokenNewResponse{
						Value: "generated-token-value",
						ID:    "new-token-id",
					}, nil)
			},
			want: &TokenGenerationResult{
				Token: "generated-token-value",
				Name:  "custom-token-name",
				Zone:  "example.com",
			},
			wantErr: false,
		},
		{
			name: "with expiration",
			req: TokenGenerationRequest{
				ServiceName: "myservice",
				ZoneName:    "example.com",
				ExpiresOn:   "2027-01-01T00:00:00Z",
			},
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{{ID: "zone123", Name: "example.com"}}, nil)
				m.EXPECT().CreateAPIToken(mock.Anything, mock.Anything).
					Return(&user.TokenNewResponse{
						Value: "generated-token-value",
						ID:    "new-token-id",
					}, nil)
			},
			want: &TokenGenerationResult{
				Token:     "generated-token-value",
				Name:      "myservice.example.com",
				Zone:      "example.com",
				ExpiresOn: "2027-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "invalid expiration format",
			req: TokenGenerationRequest{
				ServiceName: "myservice",
				ZoneName:    "example.com",
				ExpiresOn:   "not-a-date",
			},
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{{ID: "zone123", Name: "example.com"}}, nil)
			},
			wantErr: true,
		},
		{
			name: "zone not found",
			req: TokenGenerationRequest{
				ServiceName: "myservice",
				ZoneName:    "missing.com",
			},
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{}, nil)
			},
			wantErr: true,
			errType: ErrZoneNotFound,
		},
		{
			name: "token creation fails",
			req: TokenGenerationRequest{
				ServiceName: "myservice",
				ZoneName:    "example.com",
			},
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListZones(mock.Anything, mock.Anything).
					Return([]zones.Zone{{ID: "zone123", Name: "example.com"}}, nil)
				m.EXPECT().CreateAPIToken(mock.Anything, mock.Anything).
					Return(nil, errors.New("creation failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPI(t)
			tt.mockFn(mockAPI)

			got, err := GenerateToken(ctx, mockAPI, tt.req)

			if tt.wantErr {
				require.Error(t, err)

				if tt.errType != nil {
					require.ErrorIs(t, err, tt.errType)
				}

				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
