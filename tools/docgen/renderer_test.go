// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHugoRenderer(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "creates renderer",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHugoRenderer()

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_hugoRenderer_RenderCommand(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "render succeeds",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewHugoRenderer()
			require.NoError(t, err)

			doc := &CommandDoc{
				Name:  "test",
				Short: "Test",
				Use:   "test",
			}
			tmpDir := t.TempDir()

			err = r.RenderCommand(doc, tmpDir+"/test.md")

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_hugoRenderer_RenderIndex(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "render succeeds",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewHugoRenderer()
			require.NoError(t, err)

			doc := &IndexDoc{
				Title: "test",
			}
			tmpDir := t.TempDir()

			err = r.RenderIndex(doc, tmpDir+"/_index.md")

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}
