---
title: Validate
description: Check the configuration file for errors.  Verifies the file exists, is valid YAML, and contains a valid zone name.
type: docs
---

Check the configuration file for errors.

Verifies the file exists, is valid YAML, and contains a valid zone name.

### Usage

```bash
goGenerateCFToken config validate
```

### Examples

#### Validate the default config file

```bash
goGenerateCFToken config validate
```

#### Validate a specific config file

```bash
goGenerateCFToken config validate --config ./my-config.yaml
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--account-id` |  | "" | string | Cloudflare account ID stored as the generate default |
| `--token-name` |  | "" | string | Default token name stored as the generate default |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
