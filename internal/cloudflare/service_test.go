// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudflare/cloudflare-go/v7/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/cloudflare/mocks"
)

func TestListTokens(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mockFn  func(m *mocks.MockAPI)
		want    []TokenInfo
		wantErr bool
	}{
		{
			name: "successfully lists tokens",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListTokens(mock.Anything).
					Return([]user.Token{
						{ID: "token1", Name: "token-one", Status: "active"},
						{ID: "token2", Name: "token-two", Status: "disabled"},
					}, nil)
			},
			want: []TokenInfo{
				{ID: "token1", Name: "token-one", Status: "active"},
				{ID: "token2", Name: "token-two", Status: "disabled"},
			},
			wantErr: false,
		},
		{
			name: "returns empty list",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListTokens(mock.Anything).
					Return([]user.Token{}, nil)
			},
			want:    []TokenInfo{},
			wantErr: false,
		},
		{
			name: "api error propagates",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ListTokens(mock.Anything).
					Return(nil, errors.New("api failure"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPI(t)
			tt.mockFn(mockAPI)

			got, err := ListTokens(ctx, mockAPI)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRevokeToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		tokenID string
		mockFn  func(m *mocks.MockAPI)
		wantErr bool
	}{
		{
			name:    "successfully revokes token",
			tokenID: "token-id",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().RevokeToken(mock.Anything, "token-id").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "api error propagates",
			tokenID: "token-id",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().RevokeToken(mock.Anything, "token-id").
					Return(errors.New("revoke failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPI(t)
			tt.mockFn(mockAPI)

			err := RevokeToken(ctx, mockAPI, tt.tokenID)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidateCredentials(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mockFn  func(m *mocks.MockAPI)
		wantErr bool
	}{
		{
			name: "credentials are valid",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ValidateCredentials(mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "credentials are invalid",
			mockFn: func(m *mocks.MockAPI) {
				m.EXPECT().ValidateCredentials(mock.Anything).
					Return(errors.New("invalid credentials"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := mocks.NewMockAPI(t)
			tt.mockFn(mockAPI)

			err := ValidateCredentials(ctx, mockAPI)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}
