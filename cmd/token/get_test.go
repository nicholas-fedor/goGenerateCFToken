// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newGetCommand(t *testing.T) {
	tests := []struct {
		name        string
		expectedUse string
	}{
		{
			name:        "creates get command",
			expectedUse: "get <token-id>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGetCommand()
			assert.NotNil(t, got)
			assert.Equal(t, tt.expectedUse, got.Use)
			assert.NotNil(t, got.Flags().Lookup("json"))
			assert.NotNil(t, got.Flags().Lookup("token"))
			assert.NotNil(t, got.Flags().Lookup("timeout"))
			assert.Nil(t, got.Flags().Lookup("zone"))
			assert.Nil(t, got.Flags().Lookup("dry-run"))
		})
	}
}
