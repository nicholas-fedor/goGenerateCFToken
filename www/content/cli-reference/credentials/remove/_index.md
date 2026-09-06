---
title: Remove
description: Remove the stored Cloudflare API key from the OS keyring.  This does not affect any existing tokens created with the key.
type: docs
---

Remove the stored Cloudflare API key from the OS keyring.

This does not affect any existing tokens created with the key.

### Usage

```bash
goGenerateCFToken credentials remove
```

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
