---
title: Set
description: Store a Cloudflare API key for later commands.  The key is read interactively without echo, from --from-env (CF_API_TOKEN), or from --from-file. It is writte...
type: docs
---

Store a Cloudflare API key for later commands.

The key is read interactively without echo, from --from-env (CF_API_TOKEN),
or from --from-file. It is written to the OS keyring when available, otherwise
to the default credential file under the XDG config directory.

### Usage

```bash
goGenerateCFToken credentials set
```

### Examples

#### Store a new API key (interactive)

```bash
goGenerateCFToken credentials set
```

#### Store from the CF_API_TOKEN environment variable

```bash
goGenerateCFToken credentials set --from-env
```

#### Store from a file (Docker secrets)

```bash
goGenerateCFToken credentials set --from-file /run/secrets/cf_api_token
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--from-env` |  | false | bool | Read API key from CF_API_TOKEN |
| `--from-file` |  | "" | string | Read API key from the given file |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
