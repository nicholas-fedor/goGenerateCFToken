// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/cloudflare/cloudflare-go/v7/user"
	"github.com/cloudflare/cloudflare-go/v7/zones"
)

// Client wraps the Cloudflare SDK client.
type Client struct {
	*cloudflare.Client
}

// NewClient creates a Cloudflare client from an API token.
//
// Parameters:
//   - apiToken: The Cloudflare API token used for authentication.
//
// Returns:
//   - *Client: A new Cloudflare client instance.
//   - error: Non-nil if the apiToken is empty.
func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, ErrMissingCredentials
	}

	return &Client{cloudflare.NewClient(option.WithAPIToken(apiToken))}, nil
}

// API defines the Cloudflare API operations used by the service layer.
type API interface {
	ListZones(ctx context.Context, params zones.ZoneListParams) ([]zones.Zone, error)
	CreateAPIToken(ctx context.Context, params user.TokenNewParams) (*user.TokenNewResponse, error)
	ListTokens(ctx context.Context) ([]user.Token, error)
	RevokeToken(ctx context.Context, tokenID string) error
	ValidateCredentials(ctx context.Context) error
}

// Compile-time check that Client implements API.
var _ API = (*Client)(nil)

// ListZones retrieves zones matching the given parameters.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - params: Parameters to filter the zone list.
//
// Returns:
//   - []zones.Zone: The list of zones matching the parameters.
//   - error: Non-nil if the client is not initialized or the API call fails.
func (c *Client) ListZones(ctx context.Context, params zones.ZoneListParams) ([]zones.Zone, error) {
	if c.Client == nil {
		return nil, ErrClientNotInitialized
	}

	page, err := c.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListZonesFailed, err)
	}

	return page.Result, nil
}

// CreateAPIToken creates a new API token.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - params: Parameters for the new token including name and policies.
//
// Returns:
//   - *user.TokenNewResponse: The created token response containing the token value.
//   - error: Non-nil if the client is not initialized or token creation fails.
func (c *Client) CreateAPIToken(ctx context.Context, params user.TokenNewParams) (*user.TokenNewResponse, error) {
	if c.Client == nil {
		return nil, ErrClientNotInitialized
	}

	token, err := c.User.Tokens.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateTokenFailed, err)
	}

	return token, nil
}

// ListTokens retrieves all API tokens.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//
// Returns:
//   - []user.Token: The list of all API tokens.
//   - error: Non-nil if the client is not initialized or the API call fails.
func (c *Client) ListTokens(ctx context.Context) ([]user.Token, error) {
	if c.Client == nil {
		return nil, ErrClientNotInitialized
	}

	page, err := c.User.Tokens.List(ctx, user.TokenListParams{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListTokensFailed, err)
	}

	return page.Result, nil
}

// RevokeToken revokes an API token by ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - tokenID: The ID of the token to revoke.
//
// Returns:
//   - error: Non-nil if the client is not initialized or revocation fails.
func (c *Client) RevokeToken(ctx context.Context, tokenID string) error {
	if c.Client == nil {
		return ErrClientNotInitialized
	}

	_, err := c.User.Tokens.Delete(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRevokeTokenFailed, err)
	}

	return nil
}

// ValidateCredentials tests whether the API token is valid.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//
// Returns:
//   - error: Non-nil if the client is not initialized or credentials are invalid.
func (c *Client) ValidateCredentials(ctx context.Context) error {
	if c.Client == nil {
		return ErrClientNotInitialized
	}

	_, err := c.User.Tokens.List(ctx, user.TokenListParams{})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCredentialValidationFailed, err)
	}

	return nil
}
