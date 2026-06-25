---
title: Revoke
description: Revoke (delete) a Cloudflare API token by its ID.  Requires --force to confirm the destructive operation.
type: docs
---

Revoke (delete) a Cloudflare API token by its ID.

Requires --force to confirm the destructive operation.

### Usage

```bash
goGenerateCFToken token revoke <token-id>
```

### Examples

#### Revoke a token

```bash
goGenerateCFToken token revoke abc123 --force
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--force` |  | false | bool | Confirm token revocation |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
