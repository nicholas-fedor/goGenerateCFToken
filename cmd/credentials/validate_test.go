// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newValidateCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates validate command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newValidateCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "validate", got.Use)
		})
	}
}
