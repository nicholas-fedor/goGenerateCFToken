// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			name: "returns default config",
			want: Config{Zone: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Default()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errType error
	}{
		{
			name:    "empty zone returns error",
			cfg:     &Config{Zone: ""},
			wantErr: true,
			errType: ErrZoneRequired,
		},
		{
			name:    "invalid zone returns error",
			cfg:     &Config{Zone: "not a valid zone!"},
			wantErr: true,
			errType: ErrInvalidZone,
		},
		{
			name:    "valid zone passes",
			cfg:     &Config{Zone: "example.com"},
			wantErr: false,
		},
		{
			name:    "valid subdomain passes",
			cfg:     &Config{Zone: "sub.example.com"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestConfig_Format(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		format  string
		want    []byte
		wantErr bool
		errType error
	}{
		{
			name:   "yaml format",
			cfg:    &Config{Zone: "example.com"},
			format: "yaml",
			want:   []byte("zone: example.com\n"),
		},
		{
			name:   "json format",
			cfg:    &Config{Zone: "example.com"},
			format: "json",
			want:   []byte("{\n  \"zone\": \"example.com\"\n}\n"),
		},
		{
			name:   "empty format defaults to yaml",
			cfg:    &Config{Zone: "example.com"},
			format: "",
			want:   []byte("zone: example.com\n"),
		},
		{
			name:    "invalid format returns error",
			cfg:     &Config{Zone: "example.com"},
			format:  "xml",
			wantErr: true,
			errType: ErrInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cfg.Format(tt.format)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
