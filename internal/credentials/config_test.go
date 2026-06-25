// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCredentialsConfig(t *testing.T) {
	tests := []struct {
		name string
		want *CredentialsConfig
	}{
		{
			name: "creates default config",
			want: &CredentialsConfig{
				resolver: &Resolver{
					providers: []Provider{
						&EnvProvider{key: "CF_API_TOKEN"},
						&KeyringProvider{store: NewKeyringStore(KeyringService, KeyringUser)},
					},
				},
				store: NewKeyringStore(KeyringService, KeyringUser),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCredentialsConfig()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCredentialsConfig_Resolve(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *CredentialsConfig
		want    string
		wantErr bool
	}{
		{
			name: "env var set",
			setup: func() *CredentialsConfig {
				cfg := NewCredentialsConfig()
				cfg.resolver = NewResolver(&mockProvider{key: "test-key"})

				return cfg
			},
			want:    "test-key",
			wantErr: false,
		},
		{
			name: "no key found",
			setup: func() *CredentialsConfig {
				cfg := NewCredentialsConfig()
				cfg.resolver = NewResolver(&mockProvider{key: ""})

				return cfg
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup()

			got, err := cfg.Resolve(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCredentialsConfig_Set(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
		errType error
	}{
		{
			name:    "empty key returns error",
			key:     "",
			wantErr: true,
			errType: ErrEmptyAPIKey,
		},
		{
			name:    "valid key succeeds",
			key:     "test-key-12345",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewCredentialsConfig()

			err := cfg.Set(context.Background(), tt.key)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestCredentialsConfig_Delete(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "delete succeeds",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewCredentialsConfig()

			err := cfg.Delete(context.Background())

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}
