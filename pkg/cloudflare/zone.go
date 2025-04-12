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

func (c *Client) GetZoneID(ctx context.Context, zoneName string, api APIInterface) (string, error) {
	params := zones.ZoneListParams{Name: cloudflare.F(zoneName)}

	response, err := api.ListZones(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	switch len(response.Result) {
	case 0:
		return "", fmt.Errorf("%w: %s", ErrZoneNotFound, zoneName)
	case 1:
		return response.Result[0].ID, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrMultipleZonesFound, zoneName)
	}
}
