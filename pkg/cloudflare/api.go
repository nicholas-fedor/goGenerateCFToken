/*
Copyright © 2025 Nicholas Fedor <nick@nickfedor.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v4/user"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

// NewAPIClientFunc creates a new Cloudflare client, defaulting to NewClient.
var NewAPIClientFunc = NewClient

// Client wraps a Cloudflare API client for interacting with zones and tokens.
type Client struct {
	*cloudflare.Client
}

// NewClient initializes a new Cloudflare client with the provided API token.
// It returns an error if the token is missing.
func NewClient(apiToken string) (*Client, error) {
	// Initialize request options for the client.
	var opts []option.RequestOption

	// Validate and set the API token.
	switch {
	case apiToken != "":
		opts = append(opts, option.WithAPIToken(apiToken))
	default:
		return nil, ErrMissingCredentials
	}

	// Create the Cloudflare client with options.
	client := cloudflare.NewClient(opts...)

	// Return the wrapped client.
	return &Client{client}, nil
}

// ListZones retrieves a list of Cloudflare zones matching the given parameters.
// It returns an error if the client is not initialized or the API call fails.
func (c *Client) ListZones(
	ctx context.Context,
	params zones.ZoneListParams,
) (*pagination.V4PagePaginationArray[zones.Zone], error) {
	// Validate client initialization.
	if c.Client == nil {
		return nil, ErrClientNotInitialized
	}

	// Fetch zones using the provided parameters.
	zones, err := c.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListZonesFailed, err)
	}

	// Return the list of zones.
	return zones, nil
}

// CreateAPIToken generates a new Cloudflare API token with the specified parameters.
// It returns an error if the client is not initialized or the API call fails.
func (c *Client) CreateAPIToken(
	ctx context.Context,
	params user.TokenNewParams,
) (*user.TokenNewResponse, error) {
	// Validate client initialization.
	if c.Client == nil {
		return nil, ErrClientNotInitialized
	}

	// Create the API token.
	token, err := c.User.Tokens.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateTokenFailed, err)
	}

	// Return the created token response.
	return token, nil
}
