// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/securefile"
)

const (
	tokenFileName = "api_token"
	tokenFileMode = 0o600
	tokenDirMode  = 0o700
)

// FileProvider resolves an API key from a file path specified by an
// environment variable. This supports Docker Secrets and similar
// file-based credential injection patterns.
type FileProvider struct {
	envVar string
}

// NewFileProvider creates a FileProvider that reads the file path from the
// given environment variable, then reads the API key from that file.
//
// Parameters:
//   - envVar: The environment variable name containing the file path.
//
// Returns:
//   - *FileProvider: A new file-based credential provider.
func NewFileProvider(envVar string) *FileProvider {
	return &FileProvider{envVar: envVar}
}

// Resolve reads the API key from the file specified by the environment variable.
// The file content is trimmed of leading/trailing whitespace and newline characters.
//
// Parameters:
//   - ctx: Context (unused, present for interface conformance).
//
// Returns:
//   - string: The API key, or empty if not available.
//   - error: Non-nil if the file cannot be read.
func (p *FileProvider) Resolve(_ context.Context) (string, error) {
	path := os.Getenv(p.envVar)
	if path == "" {
		return "", nil
	}

	data, err := os.ReadFile(path) //nolint:gosec // G703: path comes from an explicit env var the user set
	if err != nil {
		log.Debug().Err(err).Str("path", path).Msg("failed to read API key file")

		return "", fmt.Errorf("read API key file: %w", err)
	}

	key := strings.TrimSpace(string(data))
	if key != "" {
		log.Debug().Str("source", p.envVar).Msg("resolved API key from file")
	}

	return key, nil
}

// FileStore persists an API key to a local file with owner-only permissions.
type FileStore struct {
	path string
}

// DefaultTokenFile returns the default path for a file-backed API token.
//
// Returns:
//   - string: $XDG_CONFIG_HOME/gogeneratecftoken/api_token
func DefaultTokenFile() string {
	return filepath.Join(xdg.ConfigHome, "gogeneratecftoken", tokenFileName)
}

// NewFileStore creates a FileStore that reads and writes the given path.
//
// Parameters:
//   - path: Absolute path of the credential file.
//
// Returns:
//   - *FileStore: A new file-backed credential store.
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

// Resolve reads the stored API key. A missing file is treated as unset.
//
// The file's permissions are checked and tightened first. A file that cannot be tightened
// is reported as a warning rather than an error, because the key is already exposed at a
// wider mode and refusing to read it would not reduce that exposure, so a read-only or
// otherwise unmodifiable file must not stop the command from working.
//
// Parameters:
//   - ctx: Context (unused, present for interface conformance).
//
// Returns:
//   - string: The API key, or empty if the file does not exist.
//   - error: Non-nil if the file exists but cannot be read.
func (s *FileStore) Resolve(_ context.Context) (string, error) {
	ensurePrivateMode(s.path)

	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", fmt.Errorf("read credential file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

// ensurePrivateMode tightens the credential file to owner-only permissions when it is
// accessible by group or others. A file that cannot be tightened is reported as a warning
// rather than an error, so that an unmodifiable file does not stop the credential from being
// used.
//
// Nothing is done on Windows, where Unix permission bits are not meaningful and the file
// inherits the access control list of its parent directory.
//
// Parameters:
//   - path: The credential file path.
func ensurePrivateMode(path string) {
	if runtime.GOOS == "windows" {
		return
	}

	// Lstat rather than Stat, so a symlink at the credential path is recognised as one. The
	// target of such a link is a file this tool did not write and whose permissions are the
	// user's to decide, so they are reported instead of changed. The read below still
	// follows the link, so the key itself is unaffected.
	info, err := os.Lstat(path)
	if err != nil {
		return
	}

	if info.Mode()&os.ModeSymlink != 0 {
		log.Warn().
			Str("path", path).
			Msg("credential path is a symlink - the permissions of its target are not managed here")

		return
	}

	// Only a regular file is ours to tighten. Chmod'ing a directory to 0600 would strip its
	// traversal bit and leave the path unusable, and other node types are not credentials.
	if !info.Mode().IsRegular() {
		return
	}

	mode := info.Mode().Perm()
	if mode&0o077 == 0 {
		return
	}

	log.Warn().
		Str("path", path).
		Str("permissions", fmt.Sprintf("%04o", mode)).
		Msg("credential file was accessible by group or others - tightening to 0600")

	err = os.Chmod(path, tokenFileMode)
	if err != nil {
		log.Warn().
			Err(err).
			Str("path", path).
			Msg("could not tighten credential file permissions - treat the key as exposed and rotate it")
	}
}

// Set writes the API key to the file with 0600 permissions.
//
// The key is staged in a temporary file and moved into place, and the permissions are
// enforced even when the file already exists, so a file left at a wider mode by a backup, a
// copy or an earlier version is tightened rather than written through. A symlink at the
// path is replaced rather than followed.
//
// Parameters:
//   - key: The API key to store.
//
// Returns:
//   - error: Non-nil if the directory cannot be created or the file cannot be written.
func (s *FileStore) Set(key string) error {
	dir := filepath.Dir(s.path)
	if dir != "." {
		err := os.MkdirAll(dir, tokenDirMode)
		if err != nil {
			return fmt.Errorf("create credential directory: %w", err)
		}
	}

	err := securefile.WriteFile(s.path, []byte(key+"\n"), tokenFileMode)
	if err != nil {
		return fmt.Errorf("write credential file: %w", err)
	}

	return nil
}

// Delete removes the credential file. It is idempotent when the file is absent.
//
// Returns:
//   - error: Non-nil if removal fails for a reason other than not found.
func (s *FileStore) Delete() error {
	err := os.Remove(s.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete credential file: %w", err)
	}

	return nil
}
