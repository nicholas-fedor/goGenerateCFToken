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

import "errors"

var (
	// ErrMissingCredentials indicates that the API token is not provided.
	ErrMissingCredentials = errors.New("api_token must be provided")

	// ErrClientNotInitialized indicates that the Cloudflare client is not initialized.
	ErrClientNotInitialized = errors.New("client not initialized")

	// ErrGetZoneIDFailed indicates a failure to retrieve a Cloudflare zone ID.
	ErrGetZoneIDFailed = errors.New("failed to get zone ID")

	// ErrListZonesFailed indicates a failure to list Cloudflare zones.
	ErrListZonesFailed = errors.New("failed to list zones")

	// ErrZoneNotFound indicates that no zones were found for the given name.
	ErrZoneNotFound = errors.New("no zones found")

	// ErrMultipleZonesFound indicates that multiple zones were found for the given name.
	ErrMultipleZonesFound = errors.New("multiple zones found")

	// ErrCreateTokenFailed indicates a failure to create a Cloudflare API token.
	ErrCreateTokenFailed = errors.New("failed to create API token")
)
