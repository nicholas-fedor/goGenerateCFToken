// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import "context"

// defaultConfig is the package-level credential configuration.
var defaultConfig = NewCredentialsConfig()

// ResolveAPIKey returns the Cloudflare API key.
//
// Priority: CF_API_TOKEN environment variable > OS keyring.
//
// Returns:
//   - string: The resolved API key.
//   - error: Non-nil if no API key is found in any source.
func ResolveAPIKey() (string, error) {
	return defaultConfig.Resolve(context.Background())
}

// SetAPIKey stores the API key in the OS keyring.
//
// Parameters:
//   - key: The Cloudflare API key to store.
//
// Returns:
//   - error: Non-nil if the key is empty or keyring storage fails.
func SetAPIKey(key string) error {
	return defaultConfig.Set(context.Background(), key)
}

// DeleteAPIKey removes the API key from the OS keyring.
// It returns nil if the key did not exist (idempotent).
//
// Returns:
//   - error: Non-nil if keyring deletion fails.
func DeleteAPIKey() error {
	return defaultConfig.Delete(context.Background())
}
