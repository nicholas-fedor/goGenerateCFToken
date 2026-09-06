// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newGenerateCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates generate command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGenerateCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "generate <service-name>", got.Use)
		})
	}
}

func Test_parseDuration(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{
			name:  "hours without days",
			input: "24h",
			want:  24 * time.Hour,
		},
		{
			name:  "days suffix",
			input: "7d",
			want:  7 * 24 * time.Hour,
		},
		{
			name:  "single day",
			input: "1d",
			want:  24 * time.Hour,
		},
		{
			name:  "standard go duration minutes",
			input: "30m",
			want:  30 * time.Minute,
		},
		{
			name:  "invalid empty string",
			input: "",
			want:  0,
		},
		{
			name:  "invalid single char",
			input: "x",
			want:  0,
		},
		{
			name:  "zero days rejected",
			input: "0d",
			want:  0,
		},
		{
			name:  "negative days rejected",
			input: "-1d",
			want:  0,
		},
		{
			name:  "negative go duration rejected",
			input: "-24h",
			want:  0,
		},
		{
			name:  "zero go duration rejected",
			input: "0s",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.input)
			if tt.want == 0 {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
