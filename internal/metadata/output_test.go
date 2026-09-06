// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package metadata

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintDefault(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "writes default version string",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			err := PrintDefault(writer)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			output := writer.String()
			assert.Contains(t, output, Name)
			assert.Contains(t, output, "\n")
		})
	}
}

func TestPrintVerbose(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "writes verbose output",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			err := PrintVerbose(writer)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			output := writer.String()
			assert.Contains(t, output, "name:")
			assert.Contains(t, output, "version:")
			assert.Contains(t, output, "goVersion:")
			assert.Contains(t, output, "os:")
			assert.Contains(t, output, "arch:")
		})
	}
}

func TestPrintJSON(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "writes JSON output",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			err := PrintJSON(writer)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			output := writer.String()
			assert.Contains(t, output, `"name"`)
			assert.Contains(t, output, `"version"`)
		})
	}
}
