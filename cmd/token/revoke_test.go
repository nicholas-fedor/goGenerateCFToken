// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newRevokeCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates revoke command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newRevokeCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "revoke <token-id>", got.Use)
		})
	}
}

func Test_runRevokeCmd_RequiresForce(t *testing.T) {
	err := runRevokeCmd(newRevokeCommand(), []string{"abc123"}, false)
	require.Error(t, err)
	assert.ErrorIs(t, err, errForceRequired)
}
