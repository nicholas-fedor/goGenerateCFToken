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
	// appName is the XDG application directory name.
	appName = "gogeneratecftoken"
	// currentDir is the relative path for the current directory.
	currentDir = "."
)

// DefaultDir returns the XDG config directory for this application.
//
// Returns:
//   - string: $XDG_CONFIG_HOME/gogeneratecftoken
func DefaultDir() string {
	return filepath.Join(xdg.ConfigHome, appName)
}

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
			Msg("config file is accessible by group or others - consider chmod 600")
	}
}

// Locate returns the path to the configuration file that would be used.
//
// Returns:
//   - string: The path to the first found config file, or the default XDG location.
//   - error: Always nil. Returns the default path if no config file is found.
func Locate() (string, error) {
	configDir := DefaultDir()

	candidates := []string{
		filepath.Join(configDir, configFileName),
	}

	home, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".gogeneratecftoken", configFileName),
			filepath.Join(home, ".goGenerateCFToken", configFileName),
		)
	}

	candidates = append(candidates, filepath.Join(currentDir, configFileName))

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
	log.Debug().Str("path", fallback).Msg("no config file found, using default location")

	return fallback, nil
}
