// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvProvider(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want *EnvProvider
	}{
		{
			name: "creates env provider",
			key:  "CF_API_TOKEN",
			want: &EnvProvider{key: "CF_API_TOKEN"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEnvProvider(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnvProvider_Resolve(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		envVal string
		want   string
	}{
		{
			name:   "env var set",
			key:    "TEST_ENV_VAR",
			envVal: "test-token",
			want:   "test-token",
		},
		{
			name:   "env var not set",
			key:    "UNSET_ENV_VAR",
			envVal: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.key, tt.envVal)
			p := NewEnvProvider(tt.key)

			got, err := p.Resolve(context.Background())

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewKeyringProvider(t *testing.T) {
	tests := []struct {
		name  string
		store *KeyringStore
		want  *KeyringProvider
	}{
		{
			name:  "creates keyring provider",
			store: NewKeyringStore("test-service", "test-user"),
			want:  &KeyringProvider{store: NewKeyringStore("test-service", "test-user")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKeyringProvider(tt.store)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKeyringProvider_Resolve(t *testing.T) {
	tests := []struct {
		name    string
		store   *KeyringStore
		want    string
		wantErr bool
	}{
		{
			name:    "keyring access fails",
			store:   NewKeyringStore("test-service", "test-user"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewKeyringProvider(tt.store)

			got, err := p.Resolve(context.Background())

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

func TestNewResolver(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		want      *Resolver
	}{
		{
			name:      "no providers",
			providers: nil,
			want:      &Resolver{providers: nil},
		},
		{
			name:      "with providers",
			providers: []Provider{NewEnvProvider("CF_API_TOKEN")},
			want:      &Resolver{providers: []Provider{NewEnvProvider("CF_API_TOKEN")}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewResolver(tt.providers...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResolver_Resolve(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		providers []Provider
		want      string
		wantErr   bool
	}{
		{
			name:      "no providers returns error",
			providers: nil,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "first provider returns key",
			providers: []Provider{mockProvider{key: "first-key"}},
			want:      "first-key",
			wantErr:   false,
		},
		{
			name:      "second provider used when first empty",
			providers: []Provider{mockProvider{key: ""}, mockProvider{key: "second-key"}},
			want:      "second-key",
			wantErr:   false,
		},
		{
			name:      "all providers empty returns error",
			providers: []Provider{mockProvider{key: ""}, mockProvider{key: ""}},
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewResolver(tt.providers...)

			got, err := r.Resolve(ctx)

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
