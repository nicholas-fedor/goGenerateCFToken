---
title: Set
description: Set a configuration value.  If zone is not provided, prompts for it. Currently only the zone field is supported.
type: docs
---

Set a configuration value.

If zone is not provided, prompts for it. Currently only the zone field is supported.

### Usage

```bash
goGenerateCFToken config set [zone]
```

### Examples

#### Set zone with argument

```bash
goGenerateCFToken config set example.com
```

#### Interactive prompt

```bash
goGenerateCFToken config set
```


### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
