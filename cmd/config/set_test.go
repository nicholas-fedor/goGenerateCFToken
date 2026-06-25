// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSetCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates set command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSetCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "set [zone]", got.Use)
		})
	}
}