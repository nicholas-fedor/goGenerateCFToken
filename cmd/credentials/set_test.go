// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/credentials"
)

func Test_newSetCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates set command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSetCommand()
			assert.NotNil(t, got)
			assert.Equal(t, "set", got.Use)
			assert.NotNil(t, got.Flags().Lookup("from-env"))
			assert.NotNil(t, got.Flags().Lookup("from-file"))
		})
	}
}

func Test_runSetCmd_FromFile(t *testing.T) {
	t.Setenv(credentials.EnvVarToken, "")
	t.Setenv(credentials.EnvVarTokenFile, "")

	path := filepath.Join(t.TempDir(), "token")
	require.NoError(t, os.WriteFile(path, []byte("from-file-key\n"), 0o600))

	require.NoError(t, runSetCmd(false, path))

	got, err := credentials.ResolveAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "from-file-key", got)
}

func Test_runSetCmd_FromEnvAndFileRejected(t *testing.T) {
	err := runSetCmd(true, "/tmp/token")
	require.Error(t, err)
	assert.ErrorIs(t, err, errFromEnvAndFile)
}
