// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns root command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Root()
			assert.NotNil(t, got)
			assert.Equal(t, "goGenerateCFToken", got.Use)

			gen, _, err := got.Find([]string{"generate"})
			require.NoError(t, err)
			assert.NotNil(t, gen)
			assert.NotEmpty(t, gen.Deprecated)
		})
	}
}
