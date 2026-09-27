// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package securefile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// tempPattern is the prefix given to the temporary file created before the rename. The
// leading dot keeps it out of the way of a directory listing during the brief window it
// exists.
const tempPattern = ".securefile-"

// WriteFile writes data to path, replacing any existing file in one step.
//
// The data is written to a temporary file in the same directory, which os.CreateTemp
// creates with owner-only permissions and an O_EXCL guarantee, then chmodded to perm,
// synced, and only then renamed over path. The secret is therefore never present at a
// permissive mode and never observed partially written, and because rename does not follow
// symlinks, a symlink at path is replaced rather than written through.
//
// A failed rename leaves the original file untouched, because os.Rename does not remove the
// destination first. The replacement is atomic on Unix, so a concurrent reader sees either
// the old contents or the new ones. os.Rename is explicitly documented as not an atomic
// operation on non-Unix platforms, so on Windows a reader racing the replacement may briefly
// observe neither, though never a partial write.
//
// Parameters:
//   - path: The destination path.
//   - data: The bytes to write.
//   - perm: The permissions to enforce on the file, applied even when path already exists.
//
// Returns:
//   - error: Non-nil if the destination could not be written.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, tempPattern)
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	tmpPath := tmp.Name()

	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	err = tmp.Chmod(perm)
	if err != nil {
		return fmt.Errorf("set temporary file permissions: %w", err)
	}

	_, err = tmp.Write(data)
	if err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}

	err = tmp.Sync()
	if err != nil {
		return fmt.Errorf("sync temporary file: %w", err)
	}

	err = tmp.Close()
	if err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	err = os.Rename(tmpPath, path)
	if err != nil {
		return fmt.Errorf("replace destination: %w", err)
	}

	return syncDir(dir)
}

// syncDir flushes a directory entry so that a completed rename is durable. Syncing a
// directory handle is not supported on Windows, where the write has already succeeded, so
// a failure there is not reported.
//
// Parameters:
//   - dir: The directory to sync.
//
// Returns:
//   - error: Non-nil if the directory could not be synced on a platform that supports it.
func syncDir(dir string) error {
	handle, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("open directory for sync: %w", err)
	}

	defer func() { _ = handle.Close() }()

	err = handle.Sync()
	if err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("sync directory: %w", err)
	}

	return nil
}
