// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_initVersion(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "initializes version without panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initVersion()
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns non-empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String()
			assert.NotEmpty(t, got)
		})
	}
}

func TestGetGoVersion(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns go version or unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGoVersion()
			assert.NotEmpty(t, got)
		})
	}
}

func TestFormatUTC(t *testing.T) {
	tests := []struct {
		name   string
		utcStr string
		want   string
	}{
		{
			name:   "empty string returns empty",
			utcStr: "",
			want:   "",
		},
		{
			name:   "valid UTC time converts",
			utcStr: "2026-01-01T00:00:00Z",
			want:   "2026-01-01 00:00:00 UTC",
		},
		{
			name:   "invalid time returns original",
			utcStr: "not-a-time",
			want:   "not-a-time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatUTC(tt.utcStr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns populated version info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInfo()
			assert.Equal(t, Name, got.Name)
			assert.Equal(t, Version, got.Version)
			assert.NotEmpty(t, got.GoVersion)
			assert.NotEmpty(t, got.OS)
			assert.NotEmpty(t, got.Arch)
		})
	}
}
