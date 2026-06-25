// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newGenerateCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates generate command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGenerateCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "generate <service-name>", got.Use)
		})
	}
}
