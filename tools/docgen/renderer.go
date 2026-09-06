// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

const (
	// dirPerms is the permission mode for created directories.
	dirPerms = 0o755
	// filePerms is the permission mode for created files.
	filePerms = 0o600
)

// TemplateRenderer defines the interface for rendering documentation templates.
type TemplateRenderer interface {
	// RenderCommand renders a command documentation page to the specified path.
	RenderCommand(doc *CommandDoc, path string) error
	// RenderIndex renders an index documentation page to the specified path.
	RenderIndex(doc *IndexDoc, path string) error
}

// hugoRenderer implements TemplateRenderer for Hugo-compatible markdown.
type hugoRenderer struct {
	// commandTmpl is the template for command documentation pages.
	commandTmpl *template.Template
	// indexTmpl is the template for index documentation pages.
	indexTmpl *template.Template
}

// callerFile returns the source file path of the caller using runtime.Caller.
//
// Returns:
//   - string: The caller's source file path, or empty if it cannot be determined.
func callerFile() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filename
}

// NewHugoRenderer creates a new hugoRenderer with loaded templates.
//
// Returns:
//   - TemplateRenderer: A new renderer for Hugo-compatible markdown.
//   - error: Non-nil if template parsing fails.
func NewHugoRenderer() (TemplateRenderer, error) {
	funcMap := template.FuncMap{
		"shorthand": func(s string) string {
			if s == "" {
				return ""
			}

			return "-" + s
		},
		"defaultVal": func(f FlagDoc) string {
			if f.Type == "string" && f.Default == "" {
				return `""`
			}

			return f.Default
		},
	}

	templateDir := filepath.Join(filepath.Dir(callerFile()), "templates")

	commandTmpl, err := template.New("command.tmpl").Funcs(funcMap).ParseFiles(filepath.Join(templateDir, "command.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("parse command template: %w", err)
	}

	indexTmpl, err := template.New("index.tmpl").Funcs(funcMap).ParseFiles(filepath.Join(templateDir, "index.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("parse index template: %w", err)
	}

	return &hugoRenderer{
		commandTmpl: commandTmpl,
		indexTmpl:   indexTmpl,
	}, nil
}

// RenderCommand renders a command documentation page to the specified path.
//
// Parameters:
//   - doc: The command documentation to render.
//   - path: The output file path.
//
// Returns:
//   - error: Non-nil if template execution or file write fails.
func (r *hugoRenderer) RenderCommand(doc *CommandDoc, path string) error {
	return r.renderTemplate(r.commandTmpl, doc, path)
}

// RenderIndex renders an index documentation page to the specified path.
//
// Parameters:
//   - doc: The index documentation to render.
//   - path: The output file path.
//
// Returns:
//   - error: Non-nil if template execution or file write fails.
func (r *hugoRenderer) RenderIndex(doc *IndexDoc, path string) error {
	return r.renderTemplate(r.indexTmpl, doc, path)
}

// renderTemplate executes a template with data and writes the result to a file.
//
// Parameters:
//   - tmpl: The template to execute.
//   - data: The data to pass to the template.
//   - path: The output file path.
//
// Returns:
//   - error: Non-nil if template execution or file write fails.
func (r *hugoRenderer) renderTemplate(tmpl *template.Template, data any, path string) error {
	var buf bytes.Buffer

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("execute template for %s: %w", path, err)
	}

	err = os.WriteFile(path, buf.Bytes(), filePerms)
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	return nil
}
