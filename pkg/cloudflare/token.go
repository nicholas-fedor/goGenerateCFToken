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
	"os"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/cloudflare/cloudflare-go/v4/user"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

const (
	ZoneReadPermission = "c8fed203ed3043cba015a93ad1616f1f"
	DNSWritePermission = "4755a26eedb94da69e1066d98aa820be"
)

var GenerateTokenFunc = (*Client).GenerateToken

type clientAdapter struct {
	*cloudflare.Client
}

func (c *clientAdapter) ListZones(
	ctx context.Context,
	params zones.ZoneListParams,
) (*pagination.V4PagePaginationArray[zones.Zone], error) {
	if c.Client == nil || c.Zones == nil {
		return nil, ErrZonesServiceNotInit
	}

	zones, err := c.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}

	return zones, nil
}

func (c *clientAdapter) CreateAPIToken(
	ctx context.Context,
	params user.TokenNewParams,
) (*user.TokenNewResponse, error) {
	if c.Client == nil || c.User == nil || c.User.Tokens == nil {
		return nil, ErrTokensServiceNotInit
	}

	token, err := c.User.Tokens.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create API token: %w", err)
	}

	return token, nil
}

func (c *Client) GenerateToken(ctx context.Context, serviceName, zoneName string) (string, error) {
	if c.Client == nil {
		return "", ErrClientNotInitialized
	}

	adapter := &clientAdapter{c.Client}

	zoneID, err := c.GetZoneID(ctx, zoneName, adapter)
	if err != nil {
		return "", fmt.Errorf("failed to get zone ID: %w", err)
	}

	tokenName := serviceName + "." + zoneName
	permissions := []shared.TokenPolicyPermissionGroupParam{{
		ID: cloudflare.F(ZoneReadPermission),
	}, {
		ID: cloudflare.F(DNSWritePermission),
	}}
	resources := map[string]string{
		"com.cloudflare.api.account.zone." + zoneID: "*",
	}

	policies := []shared.TokenPolicyParam{{
		Effect:           cloudflare.F(shared.TokenPolicyEffectAllow),
		PermissionGroups: cloudflare.F(permissions),
		Resources:        cloudflare.F(resources),
	}}

	params := user.TokenNewParams{
		Name:     cloudflare.F(tokenName),
		Policies: cloudflare.F(policies),
	}

	fmt.Fprintln(os.Stdout, "Generating API token:", tokenName)

	token, err := adapter.CreateAPIToken(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token.ID, nil
}
