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
	file     *FileStore
}

const (
	// KeyringService is the service name used for OS keyring storage.
	KeyringService = "gogeneratecftoken"
	// KeyringUser is the user name used for OS keyring storage.
	KeyringUser = "cloudflare"
	// EnvVarToken is the environment variable for the API token.
	EnvVarToken = "CF_API_TOKEN" //nolint:gosec // G101: environment variable name, not a credential
	// EnvVarTokenFile is the environment variable for the file path containing
	// the API token. This supports Docker Secrets and similar file-based
	// credential injection.
	EnvVarTokenFile = "CF_API_TOKEN_FILE" //nolint:gosec // G101: environment variable name, not a credential
)

// ErrEmptyAPIKey indicates an empty API key was provided.
var ErrEmptyAPIKey = errors.New("API key cannot be empty")

// NewCredentialsConfig creates a CredentialsConfig with default providers.
// The resolver chain priority is:
//  1. CF_API_TOKEN environment variable
//  2. CF_API_TOKEN_FILE file (Docker Secrets)
//  3. OS keyring (gracefully skipped when unavailable)
//  4. Default credential file under the XDG config directory
//
// Returns:
//   - *CredentialsConfig: A new config with env, file, and keyring providers.
func NewCredentialsConfig() *CredentialsConfig {
	store := NewKeyringStore(KeyringService, KeyringUser)
	file := NewFileStore(DefaultTokenFile())
	resolver := NewResolver(
		NewEnvProvider(EnvVarToken),
		NewFileProvider(EnvVarTokenFile),
		NewKeyringProvider(store),
		file,
	)

	return &CredentialsConfig{
		resolver: resolver,
		store:    store,
		file:     file,
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

	if c.store.Available() {
		return c.store.Set(key)
	}

	return c.file.Set(key)
}

// Delete removes the API key from the OS keyring and the default credential file.
//
// Parameters:
//   - ctx: Context (unused, present for symmetry).
//
// Returns:
//   - error: Non-nil if deletion fails.
func (c *CredentialsConfig) Delete(_ context.Context) error {
	return errors.Join(c.store.Delete(), c.file.Delete())
}
