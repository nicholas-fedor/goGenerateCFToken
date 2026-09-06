---
title: Get
description: Retrieve full metadata for a single Cloudflare API token by ID.  Displays token name, status, issued date, expiration, and last-used timestamp.
type: docs
---

Retrieve full metadata for a single Cloudflare API token by ID.

Displays token name, status, issued date, expiration, and last-used timestamp.

### Usage

```bash
goGenerateCFToken token get <token-id>
```

### Examples

#### Get token metadata as plain text

```bash
goGenerateCFToken token get abc123
```

#### Get token metadata as JSON

```bash
goGenerateCFToken token get abc123 --json
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--json` |  | false | bool | Output token in JSON format |
| `--timeout` |  | 30 | int | Timeout in seconds for API calls |
| `--token` | `-t` | "" | string | Cloudflare API token (prefer CF_API_TOKEN env var or keyring) |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
