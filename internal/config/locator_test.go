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

func TestLocate_IgnoresConfigInCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	here := filepath.Join(dir, configFileName)
	require.NoError(t, os.WriteFile(here, []byte("zone: \"attacker.example\"\n"), 0o600))

	// t.Chdir restores the working directory when the test ends, so the rest of the suite
	// is unaffected.
	t.Chdir(dir)

	// Guard the premise: if the file were not actually in the working directory, the
	// assertion below would pass for the wrong reason.
	assert.FileExists(t, here)

	// Point the XDG and home candidates at an empty directory so the working directory
	// holds the only config file, which is the situation the search must refuse.
	empty := t.TempDir()

	assert.Equal(t, filepath.Join(empty, configFileName), locate(empty, empty))
}

func TestLocate_PrefersXDGOverLegacyHomePaths(t *testing.T) {
	xdg := t.TempDir()
	home := t.TempDir()

	legacy := filepath.Join(home, ".gogeneratecftoken")
	require.NoError(t, os.MkdirAll(legacy, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(legacy, configFileName), []byte("zone: \"legacy.example\"\n"), 0o600))

	// With only a legacy path present it is still found, so the legacy candidates are
	// reachable rather than merely shadowed.
	assert.Equal(t, filepath.Join(legacy, configFileName), locate(xdg, home))

	// With both present the XDG path takes precedence.
	require.NoError(t, os.WriteFile(filepath.Join(xdg, configFileName), []byte("zone: \"xdg.example\"\n"), 0o600))

	assert.Equal(t, filepath.Join(xdg, configFileName), locate(xdg, home))
}

func TestLocate_FallsBackToXDGDefault(t *testing.T) {
	empty := t.TempDir()

	assert.Equal(t, filepath.Join(empty, configFileName), locate(empty, empty))
	assert.Equal(t, filepath.Join(empty, configFileName), locate(empty, ""))
}

func Test_checkFilePermissions(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "file with secure permissions",
			path: "/tmp/secure.yaml",
		},
		{
			name: "file with insecure permissions",
			path: "/tmp/insecure.yaml",
		},
		{
			name: "nonexistent file",
			path: "/tmp/nonexistent.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkFilePermissions(tt.path)
		})
	}
}
