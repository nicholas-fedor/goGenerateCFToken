// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCredentialsConfig(t *testing.T) {
	got := NewCredentialsConfig()
	require.NotNil(t, got)
	require.NotNil(t, got.store)
	require.NotNil(t, got.file)
	require.NotNil(t, got.resolver)
	assert.Len(t, got.resolver.providers, 4)
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
		{
			name: "file provider resolves",
			setup: func() *CredentialsConfig {
				cfg := NewCredentialsConfig()
				cfg.resolver = NewResolver(&mockProvider{key: ""}, &mockProvider{key: "file-key"})

				return cfg
			},
			want:    "file-key",
			wantErr: false,
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

func TestCredentialsConfig_SetFallsBackToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api_token")
	cfg := &CredentialsConfig{
		store: &KeyringStore{available: false},
		file:  NewFileStore(path),
	}

	require.NoError(t, cfg.Set(t.Context(), "file-key"))

	got, err := cfg.file.Resolve(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "file-key", got)
}

func TestCredentialsConfig_SetRejectsEmpty(t *testing.T) {
	cfg := &CredentialsConfig{
		store: &KeyringStore{available: false},
		file:  NewFileStore(filepath.Join(t.TempDir(), "api_token")),
	}

	err := cfg.Set(t.Context(), "")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyAPIKey)
}
