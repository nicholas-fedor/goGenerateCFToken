// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"context"
	"testing"

	"github.com/cloudflare/cloudflare-go/v7/user"
	"github.com/cloudflare/cloudflare-go/v7/zones"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		apiToken string
		wantErr  bool
		errType  error
	}{
		{
			name:     "empty token returns error",
			apiToken: "",
			wantErr:  true,
			errType:  ErrMissingCredentials,
		},
		{
			name:     "valid token creates client",
			apiToken: "test-token-12345",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.apiToken)

			if tt.wantErr {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.errType)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotNil(t, got.Client)
		})
	}
}

func TestClient_ListZones_NilClient(t *testing.T) {
	client := &Client{}
	_, err := client.ListZones(context.Background(), zones.ZoneListParams{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotInitialized)
}

func TestClient_CreateAPIToken_NilClient(t *testing.T) {
	client := &Client{}
	_, err := client.CreateAPIToken(context.Background(), user.TokenNewParams{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotInitialized)
}

func TestClient_ListTokens_NilClient(t *testing.T) {
	client := &Client{}
	_, err := client.ListTokens(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotInitialized)
}

func TestClient_GetToken_NilClient(t *testing.T) {
	client := &Client{}
	_, err := client.GetToken(context.Background(), "token-id")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotInitialized)
}

func TestClient_RevokeToken_NilClient(t *testing.T) {
	client := &Client{}
	err := client.RevokeToken(context.Background(), "token-id")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotInitialized)
}

func TestClient_ValidateCredentials_NilClient(t *testing.T) {
	client := &Client{}
	err := client.ValidateCredentials(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotInitialized)
}
