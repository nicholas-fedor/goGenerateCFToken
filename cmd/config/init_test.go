// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newInitCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates init command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newInitCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "init [zone]", got.Use)
		})
	}
}
