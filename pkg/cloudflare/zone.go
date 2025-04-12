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
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

// GetZoneID retrieves the ID of a Cloudflare zone by its name.
// It returns an error if no zone, multiple zones, or a listing error occurs.
func (c *Client) GetZoneID(ctx context.Context, zoneName string, api APIInterface) (string, error) {
	// Set up parameters to filter zones by name.
	params := zones.ZoneListParams{Name: cloudflare.F(zoneName)}

	// List zones matching the name.
	response, err := api.ListZones(ctx, params)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrListZonesFailed, err)
	}

	// Check the number of matching zones.
	switch len(response.Result) {
	case 0:
		// No zones found for the given name.
		return "", fmt.Errorf("%w: %s", ErrZoneNotFound, zoneName)
	case 1:
		// Return the ID of the single matching zone.
		return response.Result[0].ID, nil
	default:
		// Multiple zones found, which is ambiguous.
		return "", fmt.Errorf("%w: %s", ErrMultipleZonesFound, zoneName)
	}
}
