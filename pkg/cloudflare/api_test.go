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
	"testing"

	"github.com/cloudflare/cloudflare-go/v4"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:  "Success",
			token: "valid-token",
		},
		{
			name:    "MissingToken",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origFunc := NewAPIClientFunc

			if tt.name == "Success" {
				NewAPIClientFunc = func(_ string) (*Client, error) {
					return &Client{Client: &cloudflare.Client{}}, nil
				}
			}

			defer func() { NewAPIClientFunc = origFunc }()

			client, err := NewClient(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && client == nil {
				t.Errorf("NewClient() returned nil client")
			}
		})
	}
}
