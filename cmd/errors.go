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

package cmd

import "errors"

var (
	ErrMissingConfigAuth = errors.New("missing required credentials in config")
	ErrMissingConfigZone = errors.New("missing required zone in config")
	ErrBindAPITokenFlag  = errors.New("failed to bind api_token flag")
	ErrBindZoneFlag      = errors.New("failed to bind zone flag")
)
