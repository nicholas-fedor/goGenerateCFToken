---
title: Show
description: Display the current configuration values.  Supports YAML (default) and JSON output formats.
type: docs
---

Display the current configuration values.

Supports YAML (default) and JSON output formats.

### Usage

```bash
goGenerateCFToken config show
```

### Examples

#### Show config as YAML

```bash
goGenerateCFToken config show
```

#### Show config as JSON

```bash
goGenerateCFToken config show --format json
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--account-id` |  | "" | string | Cloudflare account ID stored as the generate default |
| `--format` | `-f` | yaml | string | Output format (yaml, json) |
| `--token-name` |  | "" | string | Default token name stored as the generate default |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
