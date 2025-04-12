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
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

var NewAPIClientFunc = NewClient

type Client struct {
	*cloudflare.Client
}

func NewClient(apiToken string) (*Client, error) {
	var opts []option.RequestOption

	switch {
	case apiToken != "":
		opts = append(opts, option.WithAPIToken(apiToken))
	default:
		return nil, ErrMissingCredentials
	}

	client := cloudflare.NewClient(opts...)

	return &Client{client}, nil
}
