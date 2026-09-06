// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"context"
	"sync"
)

// defaultConfig lazily constructs the process-wide credential configuration
// so importing the package does not probe the OS keyring.
var defaultConfig = sync.OnceValue(NewCredentialsConfig)

// ResolveAPIKey returns the Cloudflare API key.
//
// Priority: CF_API_TOKEN env var > CF_API_TOKEN_FILE file > OS keyring >
// default credential file.
//
// Returns:
//   - string: The resolved API key.
//   - error: Non-nil if no API key is found in any source.
func ResolveAPIKey() (string, error) {
	return defaultConfig().Resolve(context.Background())
}

// SetAPIKey stores the API key in the OS keyring, or the default credential
// file when no keyring backend is available.
//
// Parameters:
//   - key: The Cloudflare API key to store.
//
// Returns:
//   - error: Non-nil if the key is empty or storage fails.
func SetAPIKey(key string) error {
	return defaultConfig().Set(context.Background(), key)
}

// DeleteAPIKey removes the API key from the OS keyring and the default
// credential file. It returns nil if the key did not exist (idempotent).
//
// Returns:
//   - error: Non-nil if deletion fails.
func DeleteAPIKey() error {
	return defaultConfig().Delete(context.Background())
}
