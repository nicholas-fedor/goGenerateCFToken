// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package prompt

import (
	"io"
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

func withTerminal(t *testing.T, secret []byte, readErr error) {
	t.Helper()

	origTerm := isTerminal
	origRead := readPassword

	isTerminal = func(_ int) bool { return true }
	readPassword = func(_ int) ([]byte, error) {
		return secret, readErr
	}

	t.Cleanup(func() {
		isTerminal = origTerm
		readPassword = origRead
	})
}

func TestPrompt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "returns trimmed input",
			input: " example.com \n",
			want:  "example.com",
		},
		{
			name:    "empty input",
			input:   "\n",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "whitespace only",
			input:   "   \n",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "eof without newline",
			input:   "partial",
			wantErr: io.EOF,
		},
		{
			name:    "closed stdin",
			input:   "",
			wantErr: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withStdin(t, tt.input)

			got, err := Prompt("Zone: ")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecret(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "returns trimmed piped input",
			input: " s3cret \n",
			want:  "s3cret",
		},
		{
			name:    "empty piped input",
			input:   "\n",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "whitespace only piped input",
			input:   "   \n",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "eof without newline",
			input:   "partial",
			wantErr: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withStdin(t, tt.input)

			got, err := Secret("API key: ")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecret_Terminal(t *testing.T) {
	tests := []struct {
		name    string
		secret  []byte
		readErr error
		want    string
		wantErr error
	}{
		{
			name:   "returns trimmed hidden input",
			secret: []byte(" s3cret "),
			want:   "s3cret",
		},
		{
			name:    "empty hidden input",
			secret:  []byte(""),
			wantErr: ErrEmptyInput,
		},
		{
			name:    "whitespace only hidden input",
			secret:  []byte("   "),
			wantErr: ErrEmptyInput,
		},
		{
			name:    "read failure",
			readErr: io.ErrClosedPipe,
			wantErr: io.ErrClosedPipe,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTerminal(t, tt.secret, tt.readErr)

			got, err := Secret("API key: ")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr error
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
			name:  "trimmed y confirms",
			input: " y \n",
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
		{
			name:    "eof without newline",
			input:   "y",
			wantErr: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withStdin(t, tt.input)

			got, err := Confirm("Revoke token? (y/N) ")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.False(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
