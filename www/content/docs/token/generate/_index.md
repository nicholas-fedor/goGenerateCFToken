---
title: Generate
description: Generate a new Cloudflare API token with DNS edit permissions for the specified service.  The token name defaults to `service-name.zone`. Use --name to overr...
type: docs
---

Generate a new Cloudflare API token with DNS edit permissions for the specified service.

The token name defaults to `service-name.zone`. Use --name to override.
The zone is loaded from the config file or --zone flag.
The API token is resolved from CF_API_TOKEN env var or the OS keyring.

### Usage

```bash
goGenerateCFToken token generate <service-name>
```

### Examples

#### Generate a token using the zone from config

```bash
goGenerateCFToken token generate myapp
```

#### Generate with a custom zone

```bash
goGenerateCFToken token generate myapp --zone example.com
```

#### Generate with a custom name and expiration

```bash
goGenerateCFToken token generate myapp --name ci-deploy --expires-on 2027-01-01T00:00:00Z
```

#### Output as JSON

```bash
goGenerateCFToken token generate myapp --json
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--dry-run` |  | false | bool | Validate inputs without creating a token |
| `--expires-on` |  | "" | string | Token expiration date (RFC3339 format, e.g. 2027-01-01T00:00:00Z) |
| `--json` |  | false | bool | Output token in JSON format |
| `--name` |  | "" | string | Custom token name (default: service.zone) |
| `--output` | `-o` | "" | string | Write token to file instead of stdout |
| `--timeout` |  | 30 | int | Timeout in seconds for API calls |
| `--token` | `-t` | "" | string | Cloudflare API token (prefer CF_API_TOKEN env var or keyring) |
| `--zone` | `-z` | "" | string | Cloudflare zone name |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
