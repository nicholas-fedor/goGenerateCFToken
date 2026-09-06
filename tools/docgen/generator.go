// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// DocGenerator orchestrates documentation generation from cobra commands.
type DocGenerator struct {
	// extractor extracts documentation data from cobra commands.
	extractor DocExtractor
	// renderer renders documentation to the output format.
	renderer TemplateRenderer
}

// NewDocGenerator creates a new DocGenerator with default implementations.
//
// Returns:
//   - *DocGenerator: A new generator with default extractor and renderer.
//   - error: Non-nil if the renderer cannot be created.
func NewDocGenerator() (*DocGenerator, error) {
	extractor := NewCobraExtractor()

	renderer, err := NewHugoRenderer()
	if err != nil {
		return nil, fmt.Errorf("create renderer: %w", err)
	}

	return &DocGenerator{
		extractor: extractor,
		renderer:  renderer,
	}, nil
}

// NewDocGeneratorWithDeps creates a DocGenerator with custom dependencies.
//
// Parameters:
//   - extractor: The documentation extractor to use.
//   - renderer: The template renderer to use.
//
// Returns:
//   - *DocGenerator: A new generator with the provided dependencies.
func NewDocGeneratorWithDeps(extractor DocExtractor, renderer TemplateRenderer) *DocGenerator {
	return &DocGenerator{
		extractor: extractor,
		renderer:  renderer,
	}
}

// Generate produces documentation for the command tree starting at root.
//
// Parameters:
//   - root: The root cobra command to generate documentation for.
//   - out: The output directory path.
//
// Returns:
//   - error: Non-nil if directory creation or rendering fails.
func (g *DocGenerator) Generate(root *cobra.Command, out string) error {
	err := os.MkdirAll(out, dirPerms)
	if err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	root.DisableAutoGenTag = true
	doc := g.extractor.Extract(root, root.Name())

	return g.generateAll(doc, out)
}

// generateAll generates all documentation files from the root command doc.
//
// Parameters:
//   - doc: The root command documentation.
//   - out: The output directory path.
//
// Returns:
//   - error: Non-nil if index rendering or section generation fails.
func (g *DocGenerator) generateAll(doc *CommandDoc, out string) error {
	err := g.renderer.RenderIndex(doc.Index, filepath.Join(out, "_index.md"))
	if err != nil {
		return fmt.Errorf("render index: %w", err)
	}

	for _, section := range doc.SubCommands {
		err = g.generateSection(section, out)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateSection generates documentation for a single command section.
//
// Parameters:
//   - section: The command documentation for the section.
//   - out: The output directory path.
//
// Returns:
//   - error: Non-nil if directory creation or rendering fails.
func (g *DocGenerator) generateSection(section *CommandDoc, out string) error {
	sectionDir := filepath.Join(out, section.Name)

	err := os.MkdirAll(sectionDir, dirPerms)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", sectionDir, err)
	}

	err = g.renderer.RenderCommand(section, filepath.Join(sectionDir, "_index.md"))
	if err != nil {
		return fmt.Errorf("render command %s: %w", section.Name, err)
	}

	for _, sub := range section.SubCommands {
		err = g.generateSection(sub, sectionDir)
		if err != nil {
			return err
		}
	}

	return nil
}
