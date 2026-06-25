// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

//nolint:wrapcheck
package cloudflare

import "context"

// ListTokens retrieves all API tokens as displayable info.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - api: The API implementation to use for listing tokens.
//
// Returns:
//   - []TokenInfo: A list of token info structs with ID, name, and status.
//   - error: Non-nil if the underlying API call fails.
func ListTokens(ctx context.Context, api API) ([]TokenInfo, error) {
	tokens, err := api.ListTokens(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]TokenInfo, len(tokens))
	for i, t := range tokens {
		result[i] = TokenInfo{
			ID:     t.ID,
			Name:   t.Name,
			Status: string(t.Status),
		}
	}

	return result, nil
}

// RevokeToken revokes an API token by ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - api: The API implementation to use for revoking the token.
//   - tokenID: The ID of the token to revoke.
//
// Returns:
//   - error: Non-nil if the underlying API call fails.
func RevokeToken(ctx context.Context, api API, tokenID string) error {
	return api.RevokeToken(ctx, tokenID)
}

// ValidateCredentials tests whether the API token is valid.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - api: The API implementation to use for validation.
//
// Returns:
//   - error: Non-nil if the underlying API call fails or credentials are invalid.
func ValidateCredentials(ctx context.Context, api API) error {
	return api.ValidateCredentials(ctx)
}
