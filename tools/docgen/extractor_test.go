// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func newTitleCaser() cases.Caser {
	return cases.Title(language.English)
}

func TestNewCobraExtractor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates extractor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCobraExtractor()
			assert.NotNil(t, got)
		})
	}
}

func Test_cobraExtractor_Extract(t *testing.T) {
	tests := []struct {
		name       string
		parentPath string
		wantName   string
	}{
		{
			name:       "basic command",
			parentPath: "root",
			wantName:   "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewCobraExtractor()
			cmd := &cobra.Command{
				Use:   "test",
				Short: "Test command",
				Long:  "A test command",
				Run:   func(cmd *cobra.Command, args []string) {},
			}

			got := e.Extract(cmd, tt.parentPath)

			require.NotNil(t, got)
			assert.Equal(t, tt.wantName, got.Name)
			assert.Equal(t, "Test command", got.Short)
			assert.Equal(t, "A test command", got.Long)
		})
	}
}

func Test_cobraExtractor_buildFullPath(t *testing.T) {
	tests := []struct {
		name       string
		use        string
		parentPath string
		want       string
	}{
		{
			name:       "same as parent",
			use:        "root",
			parentPath: "root",
			want:       "root",
		},
		{
			name:       "different from parent",
			use:        "sub",
			parentPath: "root",
			want:       "root sub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cobraExtractor{}
			cmd := &cobra.Command{Use: tt.use}

			got := e.buildFullPath(cmd, tt.parentPath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cobraExtractor_buildTitle(t *testing.T) {
	tests := []struct {
		name      string
		hasSubs   bool
		wantTitle string
	}{
		{
			name:      "leaf command",
			hasSubs:   false,
			wantTitle: "Test",
		},
		{
			name:      "command group",
			hasSubs:   true,
			wantTitle: "My Group",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cobraExtractor{
				titleCaser: newTitleCaser(),
			}
			cmd := &cobra.Command{
				Use:   "test",
				Short: "My Group",
			}

			if tt.hasSubs {
				cmd = &cobra.Command{
					Use:   "test",
					Short: "My Group",
				}
				cmd.AddCommand(&cobra.Command{
					Use: "sub",
					Run: func(cmd *cobra.Command, args []string) {},
				})
			}

			got := e.buildTitle(cmd)
			assert.Equal(t, tt.wantTitle, got)
		})
	}
}

func Test_cobraExtractor_buildDescription(t *testing.T) {
	tests := []struct {
		name string
		long string
		want string
	}{
		{
			name: "short description",
			long: "Short description",
			want: "Short description",
		},
		{
			name: "long description truncates",
			long: "This is a very long description that should be truncated because it exceeds the maximum allowed length of one hundred and sixty characters for documentation purposes",
			want: "This is a very long description that should be truncated because it exceeds the maximum allowed length of one hundred and sixty characters for documentation ...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cobraExtractor{}

			got := e.buildDescription(tt.long)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cobraExtractor_buildUseLine(t *testing.T) {
	tests := []struct {
		name     string
		use      string
		fullPath string
		want     string
	}{
		{
			name:     "empty use",
			use:      "",
			fullPath: "root",
			want:     "",
		},
		{
			name:     "simple use",
			use:      "test",
			fullPath: "root",
			want:     "root",
		},
		{
			name:     "use with args",
			use:      "test <name>",
			fullPath: "root",
			want:     "root <name>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cobraExtractor{}

			got := e.buildUseLine(tt.use, tt.fullPath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cobraExtractor_buildFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags *pflag.FlagSet
		want  int
	}{
		{
			name:  "no flags",
			flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
			want:  0,
		},
		{
			name: "with flags",
			flags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("name", "", "A name flag")

				return fs
			}(),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cobraExtractor{}

			got := e.buildFlags(tt.flags)
			assert.Len(t, got, tt.want)
		})
	}
}

func Test_cobraExtractor_titleCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "lowercase",
			s:    "test",
			want: "Test",
		},
		{
			name: "already title case",
			s:    "Test",
			want: "Test",
		},
		{
			name: "multiple words",
			s:    "my command",
			want: "My Command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cobraExtractor{
				titleCaser: newTitleCaser(),
			}

			got := e.titleCase(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_firstParagraph(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "single paragraph",
			text: "First paragraph",
			want: "First paragraph",
		},
		{
			name: "multiple paragraphs",
			text: "First paragraph\n\nSecond paragraph",
			want: "First paragraph",
		},
		{
			name: "with newlines",
			text: "Line one\nLine two",
			want: "Line one Line two",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstParagraph(tt.text)
			assert.Equal(t, tt.want, got)
		})
	}
}
