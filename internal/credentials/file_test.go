// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileProvider(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		want   *FileProvider
	}{
		{
			name:   "creates file provider",
			envVar: "CF_API_TOKEN_FILE",
			want:   &FileProvider{envVar: "CF_API_TOKEN_FILE"},
		},
		{
			name:   "empty env var",
			envVar: "",
			want:   &FileProvider{envVar: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFileProvider(tt.envVar)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileProvider_Resolve(t *testing.T) {
	ctx := context.Background()

	t.Run("env var not set returns empty", func(t *testing.T) {
		p := NewFileProvider("UNSET_FILE_VAR")

		got, err := p.Resolve(ctx)

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("file exists with token", func(t *testing.T) {
		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "token")
		require.NoError(t, os.WriteFile(tokenFile, []byte("my-secret-token\n"), 0o600))

		t.Setenv("TEST_TOKEN_FILE", tokenFile)

		p := NewFileProvider("TEST_TOKEN_FILE")

		got, err := p.Resolve(ctx)

		require.NoError(t, err)
		assert.Equal(t, "my-secret-token", got)
	})

	t.Run("file content is trimmed", func(t *testing.T) {
		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "token")
		require.NoError(t, os.WriteFile(tokenFile, []byte("  my-secret-token  \n\n"), 0o600))

		t.Setenv("TEST_TOKEN_FILE_TRIM", tokenFile)

		p := NewFileProvider("TEST_TOKEN_FILE_TRIM")

		got, err := p.Resolve(ctx)

		require.NoError(t, err)
		assert.Equal(t, "my-secret-token", got)
	})

	t.Run("file does not exist returns error", func(t *testing.T) {
		t.Setenv("TEST_TOKEN_FILE_MISSING", "/nonexistent/path/token")

		p := NewFileProvider("TEST_TOKEN_FILE_MISSING")

		got, err := p.Resolve(ctx)

		require.Error(t, err)
		assert.Empty(t, got)
	})

	t.Run("empty file returns empty string", func(t *testing.T) {
		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "empty-token")
		require.NoError(t, os.WriteFile(tokenFile, []byte(""), 0o600))

		t.Setenv("TEST_TOKEN_FILE_EMPTY", tokenFile)

		p := NewFileProvider("TEST_TOKEN_FILE_EMPTY")

		got, err := p.Resolve(ctx)

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("whitespace-only file returns empty string", func(t *testing.T) {
		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "ws-token")
		require.NoError(t, os.WriteFile(tokenFile, []byte("  \n  \n"), 0o600))

		t.Setenv("TEST_TOKEN_FILE_WS", tokenFile)

		p := NewFileProvider("TEST_TOKEN_FILE_WS")

		got, err := p.Resolve(ctx)

		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestFileStore_SetGetDelete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "creds", "api_token")
	store := NewFileStore(path)

	got, err := store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Empty(t, got)

	require.NoError(t, store.Set("file-secret"))

	got, err = store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "file-secret", got)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	require.NoError(t, store.Delete())

	got, err = store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Empty(t, got)

	require.NoError(t, store.Delete())
}
