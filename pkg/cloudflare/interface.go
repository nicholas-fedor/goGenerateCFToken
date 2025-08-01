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

	"github.com/cloudflare/cloudflare-go/v5/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v5/user"
	"github.com/cloudflare/cloudflare-go/v5/zones"
)

// APIInterface defines methods for interacting with the Cloudflare API.
// It supports listing zones and creating API tokens.
type APIInterface interface {
	// ListZones retrieves a list of Cloudflare zones matching the given parameters.
	ListZones(
		ctx context.Context,
		params zones.ZoneListParams,
	) (*pagination.V4PagePaginationArray[zones.Zone], error)

	// CreateAPIToken generates a new Cloudflare API token with the specified parameters.
	CreateAPIToken(ctx context.Context, params user.TokenNewParams) (*user.TokenNewResponse, error)
}
