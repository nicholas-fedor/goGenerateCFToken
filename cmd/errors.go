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
	// ErrMissingConfigZone indicates a missing zone name in the configuration.
	ErrMissingConfigZone = errors.New("missing required zone in config")

	// ErrBindAPITokenFlag indicates a failure to bind the API token flag to the configuration.
	ErrBindAPITokenFlag = errors.New("failed to bind api_token flag")

	// ErrBindZoneFlag indicates a failure to bind the zone flag to the configuration.
	ErrBindZoneFlag = errors.New("failed to bind zone flag")
)
