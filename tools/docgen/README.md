# docgen

CLI documentation generator for goGenerateCFToken.

## Overview

docgen introspects the Cobra command tree at runtime and emits Hugo-compatible Markdown files with frontmatter. It is the single source of truth for the site's CLI reference at `www/content/cli-reference/`.

## Directory Structure

```text
tools/docgen/
    main.go           Entry point and CLI flag parsing
    models.go         Data structures for template rendering
    extractor.go      Cobra command tree introspection
    generator.go      Orchestration: output directory, recursive generation
    renderer.go       Go text/template rendering with Hugo helpers
    templates/
        command.tmpl  Single command page template
        index.tmpl    Root index page template
```

## Architecture

docgen uses three main components:

- **DocExtractor** (`extractor.go`) — Walks the `cobra.Command` tree, builds `CommandDoc` structures with titles, descriptions, flags, examples, and subcommand index data.
- **DocGenerator** (`generator.go`) — Creates the output directory, extracts the root doc, renders the root `_index.md`, then recursively renders each top-level section and its subcommands as directories of `_index.md` files.
- **TemplateRenderer** (`renderer.go`) — Loads `text/template` files and writes rendered Markdown to disk with `0o600` permissions.

Generated files include Hugo frontmatter (`title`, `description`, `type: docs`) so they are immediately consumable by the Hugo site.

## Usage

```bash
go run ./tools/docgen -out ./www/content/cli-reference
```

or via Taskfile:

```bash
task docs
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-out` | `./www/content/cli-reference` | Output directory for generated Markdown |

## Output

Running docgen produces a directory tree of `_index.md` files:

```text
www/content/cli-reference/
    _index.md                  Root index (CLI Reference)
    config/
        _index.md
        delete/_index.md
        init/_index.md
        set/_index.md
        show/_index.md
        validate/_index.md
        reset/_index.md
    credentials/
        _index.md
        remove/_index.md
        set/_index.md
        validate/_index.md
    token/
        _index.md
        generate/_index.md
        list/_index.md
        revoke/_index.md
    version/
        _index.md
```

## Inline Generation

docgen can also be triggered via `go:generate`:

```go
//go:generate go run github.com/nicholas-fedor/gogeneratecftoken/tools/docgen -out ../../www/content/cli-reference
```

## Dependencies

- `github.com/spf13/cobra` — Command tree introspection
- `github.com/spf13/pflag` — Flag metadata extraction
- `golang.org/x/text/cases` — Title-casing command names for display

## Notes

- docgen is part of the main module (no separate `go.mod`).
- Templates are loaded from relative paths `templates/command.tmpl` and `templates/index.tmpl`; run from the `tools/docgen/` working directory.
- Run `task docs` after any CLI changes that affect command descriptions, flags, or subcommands.
