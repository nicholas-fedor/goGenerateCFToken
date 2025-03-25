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
	"errors"
	"fmt"
)

// Static error variables.
var (
	ErrListZonesFailed = errors.New("failed to list zones")
	ErrMultipleZones   = errors.New("multiple zones found")
	ErrNoZonesFound    = errors.New("no zones found")
)

func GetZoneID(ctx context.Context, zone string, api APIInterface) (string, error) {
	zones, err := api.ListZones(ctx, zone)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrListZonesFailed, err)
	}

	if len(zones) > 1 {
		return "", fmt.Errorf("%w for %q", ErrMultipleZones, zone)
	}

	if len(zones) == 0 {
		return "", fmt.Errorf("%w for %q", ErrNoZonesFound, zone)
	}

	return zones[0].ID, nil
}
