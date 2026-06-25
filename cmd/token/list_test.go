// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newListCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates list command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newListCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "list", got.Use)
		})
	}
}
