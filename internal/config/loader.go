// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	yaml "go.yaml.in/yaml/v4"
)

var errConfigNotFound = errors.New("config file not found")

// Load reads configuration from file and returns a Config.
//
// Parameters:
//   - configPath: Optional explicit path to the config file. If empty, Locate() is used.
//
// Returns:
//   - *Config: The loaded configuration.
//   - error: Non-nil if the config file does not exist (when explicitly specified),
//     cannot be read, or is invalid.
func Load(configPath string) (*Config, error) {
	path, err := resolvePath(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := readAndParse(path)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("path", path).Msg("configuration loaded successfully")

	return cfg, nil
}

// resolvePath returns the explicit config path or locates the default config file.
//
// Parameters:
//   - configPath: Optional explicit path. If empty, Locate() is called.
//
// Returns:
//   - string: The resolved config file path.
//   - error: Non-nil if the explicitly specified config file does not exist.
func resolvePath(configPath string) (string, error) {
	if configPath != "" {
		log.Debug().Str("path", configPath).Msg("using explicit config path")

		info, err := os.Stat(configPath)
		if err != nil {
			return "", fmt.Errorf("%w: %s: %w", errConfigNotFound, configPath, err)
		}

		if info == nil {
			return "", fmt.Errorf("%w: %s", errConfigNotFound, configPath)
		}

		return configPath, nil
	}

	loc, err := Locate()
	if err != nil {
		return "", fmt.Errorf("locate config: %w", err)
	}

	return loc, nil
}

// readAndParse reads and parses the config file at the given path.
//
// Parameters:
//   - path: The config file path to read.
//
// Returns:
//   - *Config: The parsed configuration (defaults if file does not exist).
//   - error: Non-nil if the file exists but cannot be read, parsed, or validated.
func readAndParse(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Debug().Str("path", path).Msg("config file does not exist - using defaults")

			cfg := Default()

			return &cfg, nil
		}

		return nil, fmt.Errorf("read config file: %w", err)
	}

	log.Debug().Str("path", path).Int("size", len(data)).Msg("config file read")

	if len(data) == 0 {
		log.Debug().Str("path", path).Msg("config file is empty - using defaults")

		cfg := Default()

		return &cfg, nil
	}

	checkFilePermissions(path)

	cfg := Default()

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	log.Debug().Str("zone", cfg.Zone).Msg("config parsed")

	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}
