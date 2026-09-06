// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	yaml "go.yaml.in/yaml/v4"
)

// Config holds all application configuration.
type Config struct {
	// Zone is the Cloudflare DNS zone name for token generation.
	Zone string `json:"zone" yaml:"zone"`
	// AccountID is the optional Cloudflare account ID for multi-account setups.
	AccountID string `json:"account_id,omitempty" yaml:"account_id,omitempty"`
	// TokenName is the optional default prefix for generated token names.
	TokenName string `json:"token_name,omitempty" yaml:"token_name,omitempty"`
}

const (
	// EnvVarZone is the environment variable used as a zone fallback.
	EnvVarZone = "CF_ZONE"
)

var (
	// ErrZoneRequired indicates the zone field is empty.
	ErrZoneRequired = errors.New("zone is required")
	// ErrInvalidZone indicates the zone name does not match the expected DNS pattern.
	ErrInvalidZone = errors.New("invalid zone name")
	// ErrInvalidFormat indicates an unsupported serialization format was requested.
	ErrInvalidFormat = errors.New("unsupported format: use yaml or json")
)

// ZonePattern is the regex for validating DNS zone names.
var ZonePattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$`)

// Default returns a Config with sensible defaults.
//
// Returns:
//   - Config: A new Config with empty zone.
func Default() Config {
	return Config{
		Zone: "",
	}
}

// Validate checks the configuration for required fields and valid values.
//
// Returns:
//   - error: Non-nil if the zone is empty or invalid.
func (cfg *Config) Validate() error {
	if cfg.Zone == "" {
		return fmt.Errorf("config error: %w", ErrZoneRequired)
	}

	if !ZonePattern.MatchString(cfg.Zone) {
		return fmt.Errorf("config error: %w: %s", ErrInvalidZone, cfg.Zone)
	}

	return nil
}

// Format serializes the configuration in the given format (yaml or json).
//
// Parameters:
//   - format: The output format ("yaml", "json", or empty for YAML).
//
// Returns:
//   - []byte: The serialized configuration.
//   - error: Non-nil if marshaling fails or the format is unsupported.
func (cfg *Config) Format(format string) ([]byte, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal json: %w", err)
		}

		return append(data, '\n'), nil
	case "yaml", "":
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("marshal yaml: %w", err)
		}

		return data, nil
	default:
		return nil, fmt.Errorf("config error: %w: %s", ErrInvalidFormat, format)
	}
}
