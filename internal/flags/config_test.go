// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFlags_Bind(t *testing.T) {
	tests := []struct {
		name string
		cf   *ConfigFlags
	}{
		{
			name: "binds empty config flags",
			cf:   &ConfigFlags{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.cf.Bind(fs)
			assert.NotNil(t, fs.Lookup("account-id"))
			assert.NotNil(t, fs.Lookup("token-name"))
		})
	}
}

func TestConfigFlags_BindShow(t *testing.T) {
	tests := []struct {
		name string
		cf   *ConfigFlags
	}{
		{
			name: "binds show flags",
			cf:   &ConfigFlags{},
		},
		{
			name: "binds with pre-existing format",
			cf: &ConfigFlags{
				LogLevel: "debug",
				Format:   "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.cf.BindShow(fs)

			format, err := fs.GetString("format")
			require.NoError(t, err)
			assert.Equal(t, "yaml", format)
		})
	}
}

func TestConfigFlags_BindShowValues(t *testing.T) {
	cf := &ConfigFlags{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cf.BindShow(fs)

	err := fs.Parse([]string{"--format", "json"})
	require.NoError(t, err)

	assert.Equal(t, "json", cf.Format)
}
