---
title: Version
description: Print the application version, including commit SHA and build details.
type: docs
---

Print the application version, including commit SHA and build details.

### Usage

```bash
goGenerateCFToken version
```

### Examples

#### Print version

```bash
goGenerateCFToken version
```

#### Print detailed version info

```bash
goGenerateCFToken version --verbose
```

#### Print version as JSON

```bash
goGenerateCFToken version --json
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--json` |  | false | bool | Output version information in JSON format |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
