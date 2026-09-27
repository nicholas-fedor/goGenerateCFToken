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
// The search deliberately does not consider the current working directory. A config file
// there would be picked up silently by anyone running the tool from a freshly cloned
// repository, a scratch directory, or any location another user controls, and a generated
// token would then be minted against whatever zone and account that file names while the
// user believes it was their own. A local config is still fully supported by naming it
// explicitly with --config, which never consults this search.
//
// Returns:
//   - string: The path to the first found config file, or the default XDG location.
//   - error: Always nil. Returns the default path if no config file is found.
func Locate() (string, error) {
	return locate(DefaultDir(), homeDir()), nil
}

// homeDir returns the user's home directory, or an empty string when it cannot be
// determined, in which case the legacy home paths are skipped.
//
// Returns:
//   - string: The home directory, or an empty string.
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}

// locate returns the first candidate path that holds a config file, or the XDG default when
// none does. The directories are parameters rather than read from the environment so the
// search order can be exercised without depending on the machine running the test.
//
// Parameters:
//   - configDir: The XDG config directory for this application.
//   - home: The user's home directory, or an empty string to skip the legacy paths.
//
// Returns:
//   - string: The path to the first found config file, or filepath.Join(configDir, configFileName).
func locate(configDir, home string) string {
	candidates := []string{
		filepath.Join(configDir, configFileName),
	}

	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".gogeneratecftoken", configFileName),
			filepath.Join(home, ".goGenerateCFToken", configFileName),
		)
	}

	log.Debug().Str("config_dir", configDir).Msg("searching for config file")

	for _, path := range candidates {
		info, err := os.Stat(path)
		if err == nil && info != nil {
			log.Debug().Str("path", path).Msg("found config file")

			return path
		}

		log.Debug().Str("path", path).Msg("config file not found")
	}

	fallback := filepath.Join(configDir, configFileName)
	log.Debug().Str("path", fallback).Msg("no config file found, using default location")

	return fallback
}
