// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package securefile

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFile_CreatesFileWithRequestedMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret")

	require.NoError(t, WriteFile(path, []byte("value\n"), 0o600))

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "value\n", string(got))

	info, err := os.Stat(path)
	require.NoError(t, err)

	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
}

func TestWriteFile_TightensExistingPermissiveFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret")
	require.NoError(t, os.WriteFile(path, []byte("old\n"), 0o644))
	// os.WriteFile applies the mode subject to the process umask, so a restrictive umask
	// would leave the file owner-only and the write under test nothing to tighten.
	require.NoError(t, os.Chmod(path, 0o644))

	require.NoError(t, WriteFile(path, []byte("new\n"), 0o600))

	info, err := os.Stat(path)
	require.NoError(t, err)

	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "new\n", string(got))
}

func TestWriteFile_HonoursWiderRequestedMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret")
	require.NoError(t, os.WriteFile(path, []byte("old\n"), 0o600))

	require.NoError(t, WriteFile(path, []byte("new\n"), 0o644))

	info, err := os.Stat(path)
	require.NoError(t, err)

	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
	}
}

func TestWriteFile_ReplacesSymlinkInsteadOfFollowingIt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating a symlink requires elevated privileges on Windows")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	path := filepath.Join(dir, "secret")

	require.NoError(t, os.WriteFile(target, []byte("target-contents\n"), 0o600))
	require.NoError(t, os.Symlink(target, path))

	require.NoError(t, WriteFile(path, []byte("new\n"), 0o600))

	got, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "target-contents\n", string(got))

	got, err = os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "new\n", string(got))

	info, err := os.Lstat(path)
	require.NoError(t, err)
	assert.NotEqual(t, os.ModeSymlink, info.Mode()&os.ModeSymlink)
}

func TestWriteFile_LeavesNoTemporaryFileBehind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret")

	require.NoError(t, WriteFile(path, []byte("value\n"), 0o600))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "secret", entries[0].Name())
}

func TestWriteFile_FailsWhenParentDirectoryMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent", "secret")

	err := WriteFile(path, []byte("value\n"), 0o600)
	require.ErrorIs(t, err, os.ErrNotExist)
}
