---
title: Set
description: Prompt for a Cloudflare API key and store it securely in the OS keyring.  The key is used by other commands to authenticate with the Cloudflare API.
type: docs
---

Prompt for a Cloudflare API key and store it securely in the OS keyring.

The key is used by other commands to authenticate with the Cloudflare API.

### Usage

```bash
goGenerateCFToken credentials set
```

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
