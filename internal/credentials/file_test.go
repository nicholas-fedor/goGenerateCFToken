// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
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

	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}

	require.NoError(t, store.Delete())

	got, err = store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Empty(t, got)

	require.NoError(t, store.Delete())
}

func TestFileStore_ResolveTightensPermissiveFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits are not meaningful on Windows")
	}

	path := filepath.Join(t.TempDir(), "api_token")
	require.NoError(t, os.WriteFile(path, []byte("exposed-secret\n"), 0o644))
	require.NoError(t, os.Chmod(path, 0o644))

	store := NewFileStore(path)

	got, err := store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "exposed-secret", got)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestFileStore_ResolveLeavesNarrowFileAlone(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits are not meaningful on Windows")
	}

	path := filepath.Join(t.TempDir(), "api_token")
	require.NoError(t, os.WriteFile(path, []byte("private-secret\n"), 0o600))
	require.NoError(t, os.Chmod(path, 0o600))

	store := NewFileStore(path)

	got, err := store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "private-secret", got)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestFileStore_ResolveLeavesDirectoryPermissionsAlone(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits are not meaningful on Windows")
	}

	path := filepath.Join(t.TempDir(), "api_token")
	require.NoError(t, os.MkdirAll(path, 0o755))
	require.NoError(t, os.Chmod(path, 0o755))

	// Reading a directory as a credential fails on every platform, so both results are
	// discarded; what matters here is that the directory keeps its traversal bit.
	_, _ = NewFileStore(path).Resolve(t.Context())

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
	assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
}

func TestFileStore_ResolveLeavesSymlinkTargetPermissionsAlone(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating a symlink requires elevated privileges on Windows")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "managed-secret")
	path := filepath.Join(dir, "api_token")

	require.NoError(t, os.WriteFile(target, []byte("managed\n"), 0o644))
	require.NoError(t, os.Chmod(target, 0o644))
	require.NoError(t, os.Symlink(target, path))

	store := NewFileStore(path)

	// The key is still read through the link; only the target's mode is left alone.
	got, err := store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "managed", got)

	info, err := os.Stat(target)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
}

func TestFileStore_ResolveMissingFileIsNoOp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent")

	got, err := NewFileStore(path).Resolve(t.Context())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestFileStore_SetTightensExistingPermissiveFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "creds", "api_token")
	require.NoError(t, os.MkdirAll(filepath.Dir(path), tokenDirMode))
	require.NoError(t, os.WriteFile(path, []byte("old-secret\n"), 0o644))
	// os.WriteFile applies the mode subject to the process umask, so a restrictive umask
	// would leave the file owner-only and the write under test nothing to tighten.
	require.NoError(t, os.Chmod(path, 0o644))

	store := NewFileStore(path)
	require.NoError(t, store.Set("new-secret"))

	info, err := os.Stat(path)
	require.NoError(t, err)

	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}

	got, err := store.Resolve(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "new-secret", got)
}
