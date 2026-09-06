// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates command successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "version", got.Use)
		})
	}
}
