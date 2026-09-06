// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
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
// Parameters:
//   - ctx: Context (unused, present for interface conformance).
//
// Returns:
//   - string: The API key, or empty if the file does not exist.
//   - error: Non-nil if the file exists but cannot be read.
func (s *FileStore) Resolve(_ context.Context) (string, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", fmt.Errorf("read credential file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

// Set writes the API key to the file with 0600 permissions.
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

	err := os.WriteFile(s.path, []byte(key+"\n"), tokenFileMode)
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
