// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	got := NewCommand()
	require.NotNil(t, got)
	assert.Equal(t, "config", got.Use)

	names := make([]string, 0, len(got.Commands()))
	for _, c := range got.Commands() {
		names = append(names, c.Name())
	}

	assert.Contains(t, names, "set")
	assert.Contains(t, names, "reset")
	assert.Contains(t, names, "delete")
	assert.Contains(t, names, "init")
}
