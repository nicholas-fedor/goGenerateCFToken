// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newResetCommand(t *testing.T) {
	got := newResetCommand()
	require.NotNil(t, got)
	assert.Equal(t, "reset", got.Use)
	assert.NotNil(t, got.Flags().Lookup("yes"))
	assert.Nil(t, got.Flags().Lookup("force"))
}

func Test_runResetCmd_DeclineCancels(t *testing.T) {
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

	err = runResetCmd(newResetCommand(), false)
	require.Error(t, err)
	assert.ErrorIs(t, err, errResetCancelled)
}
