// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newListCommand(t *testing.T) {
	tests := []struct {
		name           string
		expectedUse    string
		expectedFlags  []string
		expectedGroups []string
	}{
		{
			name:          "creates list command",
			expectedUse:   "list",
			expectedFlags: []string{"output", "json", "filter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newListCommand()
			require.NotNil(t, got)
			assert.Equal(t, tt.expectedUse, got.Use)

			for _, flag := range tt.expectedFlags {
				f := got.Flags().Lookup(flag)
				assert.NotNil(t, f, "flag %q should exist", flag)
			}
		})
	}
}
