// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package cmd is the root CLI command.
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		})
	}
}
