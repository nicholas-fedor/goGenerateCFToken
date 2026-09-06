// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/flags"
)

func Test_outputToken(t *testing.T) {
	result := &cloudflare.TokenGenerationResult{
		Token: "tok-value",
		Name:  "svc.example.com",
		Zone:  "example.com",
	}

	t.Run("format none with output writes file and no stdout", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "token.txt")
		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, &flags.GenerateFlags{Format: "none", Output: out}, result)
		require.NoError(t, err)

		data, err := os.ReadFile(out)
		require.NoError(t, err)
		assert.Equal(t, "tok-value\n", string(data))
		assert.Empty(t, buf.String())
	})

	t.Run("json with output writes json file", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "token.json")
		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, &flags.GenerateFlags{JSON: true, Format: "text", Output: out}, result)
		require.NoError(t, err)

		data, err := os.ReadFile(out)
		require.NoError(t, err)

		var parsed map[string]string
		require.NoError(t, json.Unmarshal(data, &parsed))
		assert.Equal(t, "tok-value", parsed["token"])
		assert.Equal(t, "svc.example.com", parsed["name"])
		assert.Empty(t, buf.String())
	})

	t.Run("format none without output writes nothing", func(t *testing.T) {
		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, &flags.GenerateFlags{Format: "none"}, result)
		require.NoError(t, err)
		assert.Empty(t, buf.String())
	})

	t.Run("plain stdout", func(t *testing.T) {
		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, &flags.GenerateFlags{Format: "text"}, result)
		require.NoError(t, err)
		assert.Equal(t, "tok-value\n", buf.String())
	})
}
