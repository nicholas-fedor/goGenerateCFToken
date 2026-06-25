// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
)

const (
	// configFileName is the default configuration file name.
	configFileName = "config.yaml"
	// currentDir is the relative path for the current directory.
	currentDir = "."
)

// checkFilePermissions logs a warning if the config file is accessible by group or others.
//
// Parameters:
//   - path: The file path to check permissions on.
func checkFilePermissions(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	mode := info.Mode().Perm()
	if mode&0o077 != 0 {
		log.Warn().
			Str("path", path).
			Str("permissions", fmt.Sprintf("%04o", mode)).
			Msg("config file is accessible by group or others; consider chmod 600")
	}
}

// Locate returns the path to the configuration file that would be used.
//
// Returns:
//   - string: The path to the first found config file, or the default XDG location.
//   - error: Always nil; returns the default path if no config file is found.
func Locate() (string, error) {
	configDir := filepath.Join(xdg.ConfigHome, "gogeneratecftoken")

	candidates := []string{
		filepath.Join(configDir, configFileName),
		filepath.Join(currentDir, configFileName),
	}

	log.Debug().Str("config_dir", configDir).Msg("searching for config file")

	for _, path := range candidates {
		info, err := os.Stat(path)
		if err == nil && info != nil {
			log.Debug().Str("path", path).Msg("found config file")

			return path, nil
		}

		log.Debug().Str("path", path).Msg("config file not found")
	}

	fallback := filepath.Join(configDir, configFileName)
	log.Debug().Str("path", fallback).Msg("no config file found; using default location")

	return fallback, nil
}
