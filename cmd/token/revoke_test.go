// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"os"
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
			assert.NotNil(t, got.Flags().Lookup("yes"))
			assert.Nil(t, got.Flags().Lookup("force"))
		})
	}
}

func Test_runRevokeCmd_DeclineCancels(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdin
	os.Stdin = reader

	t.Cleanup(func() {
		os.Stdin = orig
		_ = reader.Close()
	})

	_, err = writer.WriteString("n\n")
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	err = runRevokeCmd(newRevokeCommand(), []string{"abc123"}, false)
	require.Error(t, err)
	assert.ErrorIs(t, err, errRevokeCancelled)
}
