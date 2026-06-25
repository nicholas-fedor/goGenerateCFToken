// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocGeneratorWithDeps(t *testing.T) {
	tests := []struct {
		name      string
		extractor DocExtractor
		renderer  TemplateRenderer
		wantErr   bool
	}{
		{
			name:      "creates with valid deps",
			extractor: NewCobraExtractor(),
			renderer:  &mockRenderer{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDocGeneratorWithDeps(tt.extractor, tt.renderer)
			assert.NotNil(t, got)
		})
	}
}

func TestDocGenerator_generateAll(t *testing.T) {
	tests := []struct {
		name    string
		doc     *CommandDoc
		wantErr bool
	}{
		{
			name:    "nil subcommands succeeds",
			doc:     &CommandDoc{Name: "root", Index: &IndexDoc{}},
			wantErr: false,
		},
		{
			name: "with subcommands",
			doc: &CommandDoc{
				Name:  "root",
				Index: &IndexDoc{},
				SubCommands: []*CommandDoc{
					{Name: "sub1", Index: &IndexDoc{}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewDocGeneratorWithDeps(&mockExtractor{}, &mockRenderer{})
			tmpDir := t.TempDir()

			err := g.generateAll(tt.doc, tmpDir)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDocGenerator_generateSection(t *testing.T) {
	tests := []struct {
		name    string
		section *CommandDoc
		wantErr bool
	}{
		{
			name:    "simple section",
			section: &CommandDoc{Name: "section", Index: &IndexDoc{}},
			wantErr: false,
		},
		{
			name: "nested subcommands",
			section: &CommandDoc{
				Name:  "section",
				Index: &IndexDoc{},
				SubCommands: []*CommandDoc{
					{Name: "sub", Index: &IndexDoc{}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewDocGeneratorWithDeps(&mockExtractor{}, &mockRenderer{})
			tmpDir := t.TempDir()

			err := g.generateSection(tt.section, tmpDir)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

type mockRenderer struct{}

func (m *mockRenderer) RenderCommand(doc *CommandDoc, path string) error {
	return nil
}

func (m *mockRenderer) RenderIndex(doc *IndexDoc, path string) error {
	return nil
}

type mockExtractor struct{}

func (m *mockExtractor) Extract(cmd *cobra.Command, parentPath string) *CommandDoc {
	return &CommandDoc{
		Name:     cmd.Name(),
		Short:    cmd.Short,
		Long:     cmd.Long,
		Use:      cmd.Use,
		FullPath: parentPath + " " + cmd.Name(),
	}
}
