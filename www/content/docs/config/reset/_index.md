---
title: Reset
description: Reset the configuration file to default values.  Requires --force to confirm the destructive operation.
type: docs
---

Reset the configuration file to default values.

Requires --force to confirm the destructive operation.

### Usage

```bash
goGenerateCFToken config reset
```

### Examples

#### Reset config to defaults

```bash
goGenerateCFToken config reset --force
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--force` |  | false | bool | Actually perform reset |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
