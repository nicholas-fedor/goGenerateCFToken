// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package prompt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withStdin(t *testing.T, input string) {
	t.Helper()

	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdin
	os.Stdin = reader

	t.Cleanup(func() {
		os.Stdin = orig

		_ = reader.Close()
	})

	_, err = writer.WriteString(input)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "lowercase y confirms",
			input: "y\n",
			want:  true,
		},
		{
			name:  "uppercase Y confirms",
			input: "Y\n",
			want:  true,
		},
		{
			name:  "n declines",
			input: "n\n",
			want:  false,
		},
		{
			name:  "empty declines",
			input: "\n",
			want:  false,
		},
		{
			name:  "yes word declines",
			input: "yes\n",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withStdin(t, tt.input)

			got, err := Confirm("Revoke token? (y/N) ")

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
