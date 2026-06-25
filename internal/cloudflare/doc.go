// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package cloudflare provides functionality to interact with the Cloudflare API
// for managing API tokens with DNS edit permissions.
//
// The package defines a Client type that wraps the Cloudflare SDK client,
// implementing the API interface for zone and token operations. The service
// layer orchestrates token generation by resolving zone IDs and constructing
// token policies with zone read and DNS write permissions.
//
// Key components:
//   - Client: Wraps the Cloudflare SDK client, implementing the API interface.
//   - API: Defines methods for zones, tokens, and credential validation.
//   - GenerateToken: Creates a token for a given service and zone.
//   - ListTokens / RevokeToken: Token lifecycle operations.
package cloudflare
