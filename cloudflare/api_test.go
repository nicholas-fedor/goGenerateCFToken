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
	"errors"
	"testing"

	"github.com/cloudflare/cloudflare-go"
)

func TestNewAPIClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		newFunc func(token string, opts ...cloudflare.Option) (*cloudflare.API, error)
		wantErr bool
	}{
		{
			name:    "Success",
			token:   "valid-token",
			newFunc: func(_ string, _ ...cloudflare.Option) (*cloudflare.API, error) { return &cloudflare.API{}, nil },
		},
		{
			name:    "Error",
			token:   "invalid-token",
			newFunc: func(_ string, _ ...cloudflare.Option) (*cloudflare.API, error) { return nil, errors.New("auth error") },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save the original function and restore it after the test
			origFunc := newWithAPITokenFunc
			defer func() { newWithAPITokenFunc = origFunc }()

			// Mock the underlying cloudflare function
			newWithAPITokenFunc = tt.newFunc

			_, err := NewAPIClient(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPIClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
