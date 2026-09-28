// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package credentials provides Cloudflare API credential management.
//
// The package supports resolving API keys from multiple sources in priority
// order:
//
//  1. CF_API_TOKEN_FILE file path (supports Docker Secrets)
//  2. CF_API_TOKEN environment variable
//  3. OS keyring (gracefully skipped when unavailable)
//  4. Default credential file ($XDG_CONFIG_HOME/gogeneratecftoken/api_token)
//
// The file-based source comes first because the variable holds a path rather
// than a credential, so preferring it keeps the secret out of the process
// environment. A deployment that sets both resolves the file.
//
// # Logging
//
// Every log call in this package records the name of the source it resolved
// from, or the path of the file involved, and never the key itself. That is what
// keeps a master API key out of the log, and it should be preserved: a new log
// call here takes a source or a path, not a resolved value.
package credentials
