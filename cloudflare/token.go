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

	"github.com/cloudflare/cloudflare-go"
)

func GenerateToken(serviceName string, zone string, api APIInterface, ctx context.Context) (string, error) {
	// Get the Zone ID from the zone name
	zoneID, err := GetZoneID(zone, api, ctx)
	if err != nil {
		return "", err
	}

	// Specify token name
	tokenName := serviceName + "." + zone

	// Output input values
	fmt.Printf("Generating API token: %s\n", tokenName)

	// Specify API token to create
	resources := make(map[string]any)
	resources["com.cloudflare.api.account.zone."+zoneID] = "*"
	tokenToCreate := cloudflare.APIToken{
		Name: tokenName,
		Policies: []cloudflare.APITokenPolicies{{
			Effect:    "allow",
			Resources: resources,
			PermissionGroups: []cloudflare.APITokenPermissionGroups{
				{
					ID:   "c8fed203ed3043cba015a93ad1616f1f",
					Name: "Zone Read",
				},
				{
					ID:   "4755a26eedb94da69e1066d98aa820be",
					Name: "DNS Write",
				},
			},
		}},
	}

	// Send the request to the Cloudflare API
	generatedToken, err := api.CreateAPIToken(ctx, tokenToCreate)
	if err != nil {
		return "", fmt.Errorf("failed to create API token: %v", err)
	}

	return generatedToken.Value, nil
}
