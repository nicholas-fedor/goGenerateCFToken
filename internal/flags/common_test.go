// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonFlags_Bind(t *testing.T) {
	tests := []struct {
		name string
		cf   *CommonFlags
	}{
		{
			name: "binds flags successfully",
			cf:   &CommonFlags{},
		},
		{
			name: "binds with pre-existing values",
			cf: &CommonFlags{
				LogLevel: "debug",
				Config:   "/path/to/config",
				Quiet:    true,
				Verbose:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.cf.Bind(fs)

			logLevel, err := fs.GetString("log-level")
			require.NoError(t, err)
			assert.Equal(t, "info", logLevel)

			config, err := fs.GetString("config")
			require.NoError(t, err)
			assert.Empty(t, config)

			quiet, err := fs.GetBool("quiet")
			require.NoError(t, err)
			assert.False(t, quiet)

			verbose, err := fs.GetBool("verbose")
			require.NoError(t, err)
			assert.False(t, verbose)
		})
	}
}

func TestCommonFlags_BindValues(t *testing.T) {
	cf := &CommonFlags{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cf.Bind(fs)

	err := fs.Parse([]string{"--log-level", "debug", "--config", "/tmp/config.yaml", "--quiet", "--verbose"})
	require.NoError(t, err)

	assert.Equal(t, "debug", cf.LogLevel)
	assert.Equal(t, "/tmp/config.yaml", cf.Config)
	assert.True(t, cf.Quiet)
	assert.True(t, cf.Verbose)
}
