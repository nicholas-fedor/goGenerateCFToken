---
title: Init
description: Initialize the configuration file.  Optionally accepts a zone as an argument. If not provided, prompts for it.
type: docs
---

Initialize the configuration file.

Optionally accepts a zone as an argument. If not provided, prompts for it.


### Usage

```bash
goGenerateCFToken config init [zone]
```

### Examples

#### Interactive initialization

```bash
goGenerateCFToken config init
```

#### Initialize with zone argument

```bash
goGenerateCFToken config init example.com
```


### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
