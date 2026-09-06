// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionFlags_Bind(t *testing.T) {
	tests := []struct {
		name string
		vf   *VersionFlags
	}{
		{
			name: "binds json flag",
			vf:   &VersionFlags{},
		},
		{
			name: "binds with pre-existing values",
			vf: &VersionFlags{
				LogLevel: "debug",
				JSON:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.vf.Bind(fs)

			json, err := fs.GetBool("json")
			require.NoError(t, err)
			assert.False(t, json)
		})
	}
}

func TestVersionFlags_BindValues(t *testing.T) {
	vf := &VersionFlags{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	vf.Bind(fs)

	err := fs.Parse([]string{"--json"})
	require.NoError(t, err)

	assert.True(t, vf.JSON)
}
