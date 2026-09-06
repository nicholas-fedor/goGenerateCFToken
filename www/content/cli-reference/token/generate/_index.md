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

#### Generate with a TTL and account ID

```bash
goGenerateCFToken token generate myapp --ttl 7d --account-id acc123
```

#### Output as JSON to stdout

```bash
goGenerateCFToken token generate myapp --json
```

#### Write token to file with no stdout

```bash
goGenerateCFToken token generate myapp --output /tmp/token.txt --format none
```


### Command Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--account-id` |  | "" | string | Cloudflare account ID for disambiguating zones across accounts |
| `--dry-run` |  | false | bool | Validate inputs without creating a token |
| `--expires-on` |  | "" | string | Token expiration date (RFC3339 format, e.g. 2027-01-01T00:00:00Z) |
| `--format` |  | text | string | Output format (text, json, none) |
| `--json` |  | false | bool | Output token in JSON format |
| `--name` |  | "" | string | Custom token name (default: service.zone) |
| `--output` | `-o` | "" | string | Write token to file instead of stdout |
| `--timeout` |  | 30 | int | Timeout in seconds for API calls |
| `--token` | `-t` | "" | string | Cloudflare API token (prefer CF_API_TOKEN env var or keyring) |
| `--ttl` |  | "" | string | Token lifetime as a human duration (e.g. 24h, 7d) |
| `--zone` | `-z` | "" | string | Cloudflare zone name |

### Global Options

| Flag | Short | Default | Type | Description |
|------|-------|---------|------|-------------|
| `--config` | `-c` | "" | string | Path to config file |
| `--log-level` | `-l` | info | string | Set logging level (debug, info, warn, error) |
| `--quiet` | `-q` | false | bool | Suppress all output except errors |
| `--verbose` | `-v` | false | bool | Enable verbose output |
