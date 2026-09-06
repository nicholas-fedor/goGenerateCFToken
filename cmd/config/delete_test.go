// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newDeleteCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates delete command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newDeleteCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "delete", got.Use)
			assert.NotNil(t, got.Flags().Lookup("yes"))
		})
	}
}
