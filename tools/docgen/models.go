// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

// CommandDoc represents a command's documentation data for template rendering.
type CommandDoc struct {
	// Name is the command name.
	Name string
	// Short is the short description.
	Short string
	// Long is the long description.
	Long string
	// Use is the usage line.
	Use string
	// Example is the example text.
	Example string
	// UseLine is the formatted usage line.
	UseLine string
	// Title is the display title.
	Title string
	// Description is the truncated description.
	Description string
	// FullPath is the full command path.
	FullPath string
	// Flags contains flag documentation.
	Flags []FlagDoc
	// Inherited contains inherited flag documentation.
	Inherited []FlagDoc
	// Examples contains parsed example documentation.
	Examples []ExampleDoc
	// SubCommands contains subcommand documentation.
	SubCommands []*CommandDoc
	// HasSubs indicates if the command has subcommands.
	HasSubs bool
	// Index contains index documentation for command groups.
	Index *IndexDoc
}

// FlagDoc represents a flag's documentation data.
type FlagDoc struct {
	// Name is the flag name.
	Name string
	// Shorthand is the single-character shorthand.
	Shorthand string
	// Default is the default value as a string.
	Default string
	// Type is the flag value type.
	Type string
	// Usage is the usage description.
	Usage string
}

// ExampleDoc represents an example block's documentation data.
type ExampleDoc struct {
	// Title is the example title.
	Title string
	// Code is the example code.
	Code string
}

// IndexDoc represents the index page's documentation data.
type IndexDoc struct {
	// Title is the index page title.
	Title string
	// Description is the index page description.
	Description string
	// Sections contains the index sections.
	Sections []SectionEntry
}

// SectionEntry represents a section in the index page.
type SectionEntry struct {
	// Name is the section name.
	Name string
	// Title is the display title.
	Title string
	// Description is the section description.
	Description string
	// HasSubs indicates if the section has subcommands.
	HasSubs bool
	// SubCommands contains the subcommand entries.
	SubCommands []SubCommandEntry
}

// SubCommandEntry represents a subcommand entry in a section.
type SubCommandEntry struct {
	// Name is the subcommand name.
	Name string
	// Description is the subcommand description.
	Description string
	// URL is the documentation URL path.
	URL string
}
