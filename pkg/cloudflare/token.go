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
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/cloudflare/cloudflare-go/v4/user"
)

// Constants defining Cloudflare permission IDs for zone read and DNS write.
const (
	// ZoneReadPermission grants read access to Cloudflare zones.
	ZoneReadPermission = "c8fed203ed3043cba015a93ad1616f1f"
	// DNSWritePermission grants write access to DNS records.
	DNSWritePermission = "4755a26eedb94da69e1066d98aa820be"
)

// GenerateTokenFunc generates a Cloudflare API token, defaulting to GenerateToken.
var GenerateTokenFunc = GenerateToken

// GenerateToken creates a new Cloudflare API token for the specified service and zone.
// It retrieves the zone ID, configures token policies, and returns the token ID.
func GenerateToken(
	ctx context.Context,
	serviceName, zoneName string,
	client *Client,
	api APIInterface,
) (string, error) {
	// Retrieve the zone ID for the given zone name.
	zoneID, err := client.GetZoneID(ctx, zoneName, api)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGetZoneIDFailed, err)
	}

	// Construct the token name from service and zone names.
	tokenName := serviceName + "." + zoneName

	// Define permissions for zone read and DNS write.
	permissions := []shared.TokenPolicyPermissionGroupParam{{
		ID: cloudflare.F(ZoneReadPermission),
	}, {
		ID: cloudflare.F(DNSWritePermission),
	}}

	// Specify resources to apply permissions to the zone.
	resources := map[string]string{
		"com.cloudflare.api.account.zone." + zoneID: "*",
	}

	// Configure token policy to allow the specified permissions and resources.
	policies := []shared.TokenPolicyParam{{
		Effect:           cloudflare.F(shared.TokenPolicyEffectAllow),
		PermissionGroups: cloudflare.F(permissions),
		Resources:        cloudflare.F(resources),
	}}

	// Set up parameters for creating the new token.
	params := user.TokenNewParams{
		Name:     cloudflare.F(tokenName),
		Policies: cloudflare.F(policies),
	}

	// Log token generation intent.
	fmt.Fprintln(os.Stdout, "Generating API token:", tokenName)

	// Create the API token.
	token, err := api.CreateAPIToken(ctx, params)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCreateTokenFailed, err)
	}

	// Return the generated token’s Value.
	return token.Value, nil
}
