// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete_CustomPathDoesNotRemoveParent(t *testing.T) {
	dir := t.TempDir()
	sibling := filepath.Join(dir, "keep.txt")
	require.NoError(t, os.WriteFile(sibling, []byte("keep"), 0o600))

	cfgPath := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte("zone: example.com\n"), 0o600))

	require.NoError(t, Delete(cfgPath))

	_, err := os.Stat(cfgPath)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(sibling)
	require.NoError(t, err)

	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestWriteAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	cfg := Config{Zone: "example.com", AccountID: "acc", TokenName: "ci"}

	require.NoError(t, Write(path, cfg))

	got, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, cfg, *got)
}
