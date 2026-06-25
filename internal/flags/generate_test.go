// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFlags_Bind(t *testing.T) {
	tests := []struct {
		name string
		gf   *GenerateFlags
	}{
		{
			name: "binds all flags",
			gf:   &GenerateFlags{},
		},
		{
			name: "binds with pre-existing values",
			gf: &GenerateFlags{
				CommonFlags: CommonFlags{LogLevel: "debug"},
				Token:       "test-token",
				Zone:        "example.com",
				Output:      "/tmp/output",
				JSON:        true,
				Name:        "my-token",
				DryRun:      true,
				Timeout:     60,
				ExpiresOn:   "2027-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.gf.Bind(fs)

			token, err := fs.GetString("token")
			require.NoError(t, err)
			assert.Empty(t, token)

			zone, err := fs.GetString("zone")
			require.NoError(t, err)
			assert.Empty(t, zone)

			output, err := fs.GetString("output")
			require.NoError(t, err)
			assert.Empty(t, output)

			json, err := fs.GetBool("json")
			require.NoError(t, err)
			assert.False(t, json)

			name, err := fs.GetString("name")
			require.NoError(t, err)
			assert.Empty(t, name)

			dryRun, err := fs.GetBool("dry-run")
			require.NoError(t, err)
			assert.False(t, dryRun)

			timeout, err := fs.GetInt("timeout")
			require.NoError(t, err)
			assert.Equal(t, DefaultTimeout, timeout)

			expiresOn, err := fs.GetString("expires-on")
			require.NoError(t, err)
			assert.Empty(t, expiresOn)
		})
	}
}

func TestGenerateFlags_BindValues(t *testing.T) {
	gf := &GenerateFlags{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	gf.Bind(fs)

	err := fs.Parse([]string{
		"--token", "my-token",
		"--zone", "example.com",
		"--output", "/tmp/out.txt",
		"--json",
		"--name", "custom-name",
		"--dry-run",
		"--timeout", "45",
		"--expires-on", "2027-06-01T00:00:00Z",
	})
	require.NoError(t, err)

	assert.Equal(t, "my-token", gf.Token)
	assert.Equal(t, "example.com", gf.Zone)
	assert.Equal(t, "/tmp/out.txt", gf.Output)
	assert.True(t, gf.JSON)
	assert.Equal(t, "custom-name", gf.Name)
	assert.True(t, gf.DryRun)
	assert.Equal(t, 45, gf.Timeout)
	assert.Equal(t, "2027-06-01T00:00:00Z", gf.ExpiresOn)
}
