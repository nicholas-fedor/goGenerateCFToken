---
title: Set
description: Set configuration values.  Zone may be given as an argument or entered interactively. Use --account-id and --token-name to store generate defaults.
type: docs
---

Set configuration values.

Zone may be given as an argument or entered interactively.
Use --account-id and --token-name to store generate defaults.

### Usage

```bash
goGenerateCFToken config set [zone]
```

### Examples

#### Set zone with argument

```bash
goGenerateCFToken config set example.com
```

#### Set account ID without changing zone

```bash
goGenerateCFToken config set --account-id acc123
```

#### Interactive prompt

```bash
goGenerateCFToken config set
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
