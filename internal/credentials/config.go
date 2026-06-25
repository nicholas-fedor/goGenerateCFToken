// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"errors"
)

// CredentialsConfig provides a unified API for managing Cloudflare credentials.
type CredentialsConfig struct {
	resolver *Resolver
	store    *KeyringStore
}

const (
	// KeyringService is the service name used for OS keyring storage.
	KeyringService = "gogeneratecftoken"
	// KeyringUser is the user name used for OS keyring storage.
	KeyringUser = "cloudflare"
)

// ErrEmptyAPIKey indicates an empty API key was provided.
var ErrEmptyAPIKey = errors.New("API key cannot be empty")

// NewCredentialsConfig creates a CredentialsConfig with default providers.
//
// Returns:
//   - *CredentialsConfig: A new config with env and keyring providers.
func NewCredentialsConfig() *CredentialsConfig {
	store := NewKeyringStore(KeyringService, KeyringUser)
	resolver := NewResolver(
		NewEnvProvider("CF_API_TOKEN"),
		NewKeyringProvider(store),
	)

	return &CredentialsConfig{
		resolver: resolver,
		store:    store,
	}
}

// Resolve returns the Cloudflare API key from the first available source.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//
// Returns:
//   - string: The resolved API key.
//   - error: Non-nil if no API key is found.
func (c *CredentialsConfig) Resolve(ctx context.Context) (string, error) {
	return c.resolver.Resolve(ctx)
}

// Set stores the API key in the OS keyring.
//
// Parameters:
//   - ctx: Context (unused, present for symmetry).
//   - key: The API key to store.
//
// Returns:
//   - error: Non-nil if the key is empty or storage fails.
func (c *CredentialsConfig) Set(_ context.Context, key string) error {
	if key == "" {
		return ErrEmptyAPIKey
	}

	return c.store.Set(key)
}

// Delete removes the API key from the OS keyring.
//
// Parameters:
//   - ctx: Context (unused, present for symmetry).
//
// Returns:
//   - error: Non-nil if deletion fails.
func (c *CredentialsConfig) Delete(_ context.Context) error {
	return c.store.Delete()
}
