---
title: Delete
description: Delete the configuration file and its parent directory.  Prompts for confirmation (y/N) unless --yes or -y is provided.
type: docs
---

Delete the configuration file and its parent directory.

Prompts for confirmation (y/N) unless --yes or -y is provided.

### Usage

```bash
goGenerateCFToken config delete
```

### Examples

#### Delete config (with confirmation prompt)

```bash
goGenerateCFToken config delete
```

#### Delete config without confirmation

```bash
goGenerateCFToken config delete --yes
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
