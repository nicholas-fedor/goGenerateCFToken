// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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

	// assertSymlinkReplaced points gflags.Output at a symlink and asserts the link is
	// replaced rather than written through, so the file behind it is left untouched.
	assertSymlinkReplaced := func(t *testing.T, gflags *flags.GenerateFlags) {
		t.Helper()

		dir := t.TempDir()
		target := filepath.Join(dir, "target")
		out := filepath.Join(dir, "token.txt")

		require.NoError(t, os.WriteFile(target, []byte("target-contents\n"), 0o600))
		require.NoError(t, os.Symlink(target, out))

		gflags.Output = out

		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, gflags, result)
		require.NoError(t, err)
		assert.Empty(t, buf.String())

		data, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, "target-contents\n", string(data))

		data, err = os.ReadFile(out)
		require.NoError(t, err)
		assert.NotEqual(t, "target-contents\n", string(data))
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

	t.Run("plain output over a permissive file tightens permissions", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "token.txt")
		// os.WriteFile applies the mode subject to the process umask, so set it explicitly
		// to guarantee the file really is permissive before the write under test.
		require.NoError(t, os.WriteFile(out, []byte("stale\n"), 0o644))
		require.NoError(t, os.Chmod(out, 0o644))

		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, &flags.GenerateFlags{Format: "none", Output: out}, result)
		require.NoError(t, err)

		data, err := os.ReadFile(out)
		require.NoError(t, err)
		assert.Equal(t, "tok-value\n", string(data))

		info, err := os.Stat(out)
		require.NoError(t, err)

		if runtime.GOOS != "windows" {
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
		}
	})

	t.Run("json output over a permissive file tightens permissions", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "token.json")
		require.NoError(t, os.WriteFile(out, []byte("{}"), 0o644))
		require.NoError(t, os.Chmod(out, 0o644))

		cmd := &cobra.Command{}

		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := outputToken(cmd, &flags.GenerateFlags{JSON: true, Format: "text", Output: out}, result)
		require.NoError(t, err)

		info, err := os.Stat(out)
		require.NoError(t, err)

		if runtime.GOOS != "windows" {
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
		}
	})

	t.Run("plain output replaces a symlink instead of writing through it", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("creating a symlink requires elevated privileges on Windows")
		}

		assertSymlinkReplaced(t, &flags.GenerateFlags{Format: "none"})
	})

	t.Run("json output replaces a symlink instead of writing through it", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("creating a symlink requires elevated privileges on Windows")
		}

		assertSymlinkReplaced(t, &flags.GenerateFlags{JSON: true, Format: "text"})
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
