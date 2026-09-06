// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package credentials provides Cloudflare API credential management.
//
// The package supports resolving API keys from multiple sources in priority
// order:
//
//  1. CF_API_TOKEN environment variable
//  2. CF_API_TOKEN_FILE file path (supports Docker Secrets)
//  3. OS keyring (gracefully skipped when unavailable)
//  4. Default credential file ($XDG_CONFIG_HOME/gogeneratecftoken/api_token)
package credentials
