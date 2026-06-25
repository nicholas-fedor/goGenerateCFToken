// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import "errors"

// ErrMissingCredentials indicates no API token was provided.
var ErrMissingCredentials = errors.New("api_token must be provided")

// ErrClientNotInitialized indicates the client is nil.
var ErrClientNotInitialized = errors.New("client not initialized")

// ErrZoneNotFound indicates no zone matched the given name.
var ErrZoneNotFound = errors.New("no zone found")

// ErrMultipleZonesFound indicates multiple zones matched the given name.
var ErrMultipleZonesFound = errors.New("multiple zones found")

// ErrListZonesFailed indicates a failure listing zones.
var ErrListZonesFailed = errors.New("failed to list zones")

// ErrCreateTokenFailed indicates a failure creating an API token.
var ErrCreateTokenFailed = errors.New("failed to create API token")

// ErrListTokensFailed indicates a failure listing API tokens.
var ErrListTokensFailed = errors.New("failed to list tokens")

// ErrRevokeTokenFailed indicates a failure revoking an API token.
var ErrRevokeTokenFailed = errors.New("failed to revoke token")

// ErrCredentialValidationFailed indicates credential validation failed.
var ErrCredentialValidationFailed = errors.New("credential validation failed")

// ErrInvalidServiceName indicates an invalid service name.
var ErrInvalidServiceName = errors.New("invalid service name: must match ^[a-zA-Z0-9_-]+$")
