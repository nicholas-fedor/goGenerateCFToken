// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newRemoveCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates remove command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newRemoveCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "remove", got.Use)
		})
	}
}
