// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newResetCommand(t *testing.T) {
	got := newResetCommand()
	require.NotNil(t, got)
	assert.Equal(t, "reset", got.Use)
	assert.NotNil(t, got.Flags().Lookup("force"))
}

func Test_runResetCmd_RequiresForce(t *testing.T) {
	err := runResetCmd(newResetCommand(), false)
	require.Error(t, err)
	assert.ErrorIs(t, err, errResetForceRequired)
}
