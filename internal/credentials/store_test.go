// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyringStore(t *testing.T) {
	tests := []struct {
		name    string
		service string
		user    string
		want    *KeyringStore
	}{
		{
			name:    "creates keyring store",
			service: "test-service",
			user:    "test-user",
			want:    &KeyringStore{service: "test-service", user: "test-user"},
		},
		{
			name:    "empty strings",
			service: "",
			user:    "",
			want:    &KeyringStore{service: "", user: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKeyringStore(tt.service, tt.user)
			assert.Equal(t, tt.want, got)
		})
	}
}
