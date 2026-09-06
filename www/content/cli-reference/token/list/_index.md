---
title: List
description: List all Cloudflare API tokens associated with the configured credentials.  Displays token ID, name, and status for each token.
type: docs
---

List all Cloudflare API tokens associated with the configured credentials.

Displays token ID, name, and status for each token.

### Usage

```bash
goGenerateCFToken token list
```

### Examples

#### List all tokens as plain text

```bash
goGenerateCFToken token list
```

#### List all tokens as JSON

```bash
goGenerateCFToken token list --json
```

#### Filter by status

```bash
goGenerateCFToken token list --filter active
```

#### Write to file

```bash
goGenerateCFToken token list --output tokens.txt
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--filter` |  | "" | string | Filter tokens by status (e.g., active, disabled) |
| `--json` |  | false | bool | Output in JSON format |
| `--output` | `-o` | "" | string | Write output to file instead of stdout |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
