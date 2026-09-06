// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	keyring "github.com/zalando/go-keyring"
)

// Provider resolves an API key from a specific source.
type Provider interface {
	// Resolve returns the API key from this provider.
	Resolve(ctx context.Context) (string, error)
}

// EnvProvider resolves an API key from an environment variable.
type EnvProvider struct {
	key string
}

// ErrNoAPIKey indicates no API key was found in any source.
var ErrNoAPIKey = errors.New("no API key found")

// NewEnvProvider creates an EnvProvider for the given environment variable name.
//
// Parameters:
//   - key: The environment variable name.
//
// Returns:
//   - *EnvProvider: A new environment variable provider.
func NewEnvProvider(key string) *EnvProvider {
	return &EnvProvider{key: key}
}

// Resolve returns the value of the environment variable if set.
//
// Parameters:
//   - ctx: Context (unused, present for interface conformance).
//
// Returns:
//   - string: Empty if not set, otherwise the value.
//   - error: Always nil.
func (p *EnvProvider) Resolve(_ context.Context) (string, error) {
	v := os.Getenv(p.key)
	if v != "" {
		log.Debug().Str("source", p.key).Msg("resolved API key from environment variable")
	}

	return v, nil
}

// KeyringProvider resolves an API key from the OS keyring.
// When the keyring is unavailable (no backend at runtime), Resolve returns
// an empty string with no error so the Resolver chain can continue to the
// next provider.
type KeyringProvider struct {
	store *KeyringStore
}

// NewKeyringProvider creates a KeyringProvider using the given store.
//
// Parameters:
//   - store: The keyring store to read from.
//
// Returns:
//   - *KeyringProvider: A new keyring provider.
func NewKeyringProvider(store *KeyringStore) *KeyringProvider {
	return &KeyringProvider{store: store}
}

// Resolve retrieves the API key from the OS keyring.
// When the keyring is unavailable, returns ("", nil) so the resolver
// chain gracefully falls through to the next provider.
//
// Parameters:
//   - ctx: Context (unused, present for interface conformance).
//
// Returns:
//   - string: The API key, or empty if unavailable.
//   - error: Non-nil only for unexpected keyring failures (not "not found"
//     or "unsupported").
func (p *KeyringProvider) Resolve(_ context.Context) (string, error) {
	if !p.store.Available() {
		log.Debug().Msg("OS keyring unavailable - skipping keyring provider")

		return "", nil
	}

	key, err := p.store.Get()
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			log.Debug().Msg("API key not found in OS keyring")

			return "", nil
		}

		return "", fmt.Errorf("read OS keyring: %w", err)
	}

	return key, nil
}

// Resolver attempts to resolve an API key from multiple providers in order.
type Resolver struct {
	providers []Provider
}

// NewResolver creates a Resolver with the given providers checked in order.
//
// Parameters:
//   - providers: Providers to check, in priority order.
//
// Returns:
//   - *Resolver: A new resolver instance.
func NewResolver(providers ...Provider) *Resolver {
	return &Resolver{providers: providers}
}

// Resolve returns the first non-empty key from the providers.
//
// Parameters:
//   - ctx: Context passed to each provider.
//
// Returns:
//   - string: The resolved API key.
//   - error: ErrNoAPIKey if no provider returns a key.
func (r *Resolver) Resolve(ctx context.Context) (string, error) {
	for _, p := range r.providers {
		key, err := p.Resolve(ctx)
		if err != nil {
			return "", fmt.Errorf("resolve API key: %w", err)
		}

		if key != "" {
			return key, nil
		}
	}

	log.Debug().Msg("API key not found in any source")

	return "", fmt.Errorf("%w: set CF_API_TOKEN or CF_API_TOKEN_FILE env var, or use 'credentials set'", ErrNoAPIKey)
}
