---
title: Revoke
description: Revoke (delete) a Cloudflare API token by its ID.  Prompts for confirmation (y/N) unless --yes or -y is provided.
type: docs
---

Revoke (delete) a Cloudflare API token by its ID.

Prompts for confirmation (y/N) unless --yes or -y is provided.

### Usage

```bash
goGenerateCFToken token revoke <token-id>
```

### Examples

#### Revoke a token (with confirmation prompt)

```bash
goGenerateCFToken token revoke abc123
```

#### Revoke a token without confirmation

```bash
goGenerateCFToken token revoke abc123 --yes
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--yes` | `-y` | false | bool | Bypass confirmation prompt |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
