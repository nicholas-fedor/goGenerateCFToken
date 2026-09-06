// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns non-nil logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Logger()
			assert.NotNil(t, got)
		})
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level zerolog.Level
	}{
		{
			name:  "debug level",
			level: zerolog.DebugLevel,
		},
		{
			name:  "info level",
			level: zerolog.InfoLevel,
		},
		{
			name:  "warn level",
			level: zerolog.WarnLevel,
		},
		{
			name:  "error level",
			level: zerolog.ErrorLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			assert.Equal(t, tt.level, zerolog.GlobalLevel())
		})
	}
}

func TestSetup_UsesRequestedLevel(t *testing.T) {
	t.Cleanup(func() {
		SetQuiet(false)
		Setup(zerolog.InfoLevel)
	})

	Setup(zerolog.ErrorLevel)
	assert.Equal(t, zerolog.ErrorLevel, zerolog.GlobalLevel())

	Setup(zerolog.WarnLevel)
	assert.Equal(t, zerolog.WarnLevel, zerolog.GlobalLevel())
}

func TestPrintf_HonorsQuiet(t *testing.T) {
	t.Cleanup(func() {
		SetQuiet(false)
	})

	SetQuiet(true)
	Printf("should-not-panic %s", "ok")
	Println("should-not-panic")
}
