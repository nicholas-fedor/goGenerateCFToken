---
title: Validate
description: Test the configured Cloudflare API credentials against the Cloudflare API.  Resolves the API key from CF_API_TOKEN env var or the OS keyring and verifies it ...
type: docs
---

Test the configured Cloudflare API credentials against the Cloudflare API.

Resolves the API key from CF_API_TOKEN env var or the OS keyring and verifies
it can authenticate successfully.

### Usage

```bash
goGenerateCFToken credentials validate
```

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
