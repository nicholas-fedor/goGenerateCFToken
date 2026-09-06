// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	yaml "go.yaml.in/yaml/v4"
)

const (
	// defaultFileMode is the secure file permission mode for config files (owner read/write only).
	defaultFileMode = 0o600
	// defaultDirMode is the permission mode for newly created config directories.
	defaultDirMode = 0o700
)

// EnsureFile creates the config file at path if it does not exist.
//
// Parameters:
//   - path: The file path to ensure exists.
//
// Returns:
//   - error: Non-nil if the file cannot be created or stat fails.
func EnsureFile(path string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		err := os.MkdirAll(dir, defaultDirMode)
		if err != nil {
			return fmt.Errorf("create config directory: %w", err)
		}
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug().Str("path", path).Msg("creating new config file")

			werr := os.WriteFile(path, []byte{}, defaultFileMode)
			if werr != nil {
				return fmt.Errorf("create config file: %w", werr)
			}
		} else {
			return fmt.Errorf("stat config file: %w", err)
		}
	} else {
		log.Debug().Str("path", path).Msg("config file already exists")
	}

	return nil
}

// Write writes the given Config to the specified path with secure permissions.
//
// Parameters:
//   - path: The file path to write to.
//   - cfg: The configuration to serialize and write.
//
// Returns:
//   - error: Non-nil if YAML marshaling or file write fails.
func Write(path string, cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	log.Debug().Str("path", path).Int("size", len(data)).Msg("writing config file")

	err = os.WriteFile(path, data, defaultFileMode)
	if err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

// Delete removes the configuration file and its parent directory.
//
// Parameters:
//   - path: The file path to delete.
//
// Returns:
//   - error: Non-nil if the file or directory cannot be removed.
func Delete(path string) error {
	dir := filepath.Dir(path)

	log.Debug().Str("path", path).Msg("deleting config file")

	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete config file: %w", err)
	}

	if filepath.Clean(dir) != filepath.Clean(DefaultDir()) {
		return nil
	}

	log.Debug().Str("dir", dir).Msg("deleting config directory")

	err = os.Remove(dir)
	if err != nil && !os.IsNotExist(err) && !isDirNotEmpty(err) {
		return fmt.Errorf("delete config directory: %w", err)
	}

	return nil
}

func isDirNotEmpty(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not empty")
}

// WriteDefaults writes the default configuration to the given path.
//
// Parameters:
//   - path: The file path to write the default configuration to.
//
// Returns:
//   - error: Non-nil if the default config cannot be written.
func WriteDefaults(path string) error {
	log.Debug().Str("path", path).Msg("writing default configuration")

	return Write(path, Default())
}
