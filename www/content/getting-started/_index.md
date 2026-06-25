---
title: Getting Started
type: docs
---

This guide walks you through installing goGenerateCFToken, configuring your credentials, and generating your first Cloudflare API token.

## Prerequisites

- [Go](https://go.dev/) 1.26+ (for installing from source)
- A [Cloudflare](https://www.cloudflare.com/) account with at least one zone
- A Cloudflare API token with the necessary permissions

## Installation

### From Source

```bash
go install github.com/nicholas-fedor/gogeneratecftoken@latest
```

### Pre-built Binary

Download the latest release from the [GitHub releases page](https://github.com/nicholas-fedor/gogeneratecftoken/releases).

## Initial Setup

### 1. Initialize Configuration

Run the interactive wizard to set up your configuration and store your Cloudflare API key:

```bash
goGenerateCFToken config init
```

This will prompt you for:
- **Cloudflare zone name** — e.g., `example.com`
- **Cloudflare API key** — stored securely in your OS keyring

The configuration file is written to `~/.config/gogeneratecftoken/config.yaml`.

### 2. Verify Credentials

Test that your credentials are valid:

```bash
goGenerateCFToken credentials validate
```

## Generating a Token

Generate a Cloudflare API token for a service:

```bash
goGenerateCFToken token generate myapp
```

This creates a token named `myapp.example.com` with:
- Zone read permissions
- DNS write permissions

The token is printed to stdout by default. Use `--output` to write it to a file:

```bash
goGenerateCFToken token generate myapp --output ./token.txt
```

### Custom Token Name

Override the default naming convention:

```bash
goGenerateCFToken token generate myapp --name ci-deploy
```

### Token Expiration

Set an expiration date for the token:

```bash
goGenerateCFToken token generate myapp --expires-on 2027-01-01T00:00:00Z
```

### JSON Output

Output the token in JSON format for use in scripts:

```bash
goGenerateCFToken token generate myapp --json
```

## Managing Tokens

### List Tokens

View all tokens associated with your Cloudflare account:

```bash
goGenerateCFToken token list
```

### Revoke a Token

Revoke a token by its ID:

```bash
goGenerateCFToken token revoke abc123 --force
```

## Managing Credentials

### Store a New API Key

```bash
goGenerateCFToken credentials set
```

### Remove Stored Credentials

```bash
goGenerateCFToken credentials remove
```

## Configuration

Configuration is loaded from an XDG-compliant path:

```
~/.config/gogeneratecftoken/config.yaml
```

Override the path with `--config`:

```bash
goGenerateCFToken --config ./custom-config.yaml token generate myapp
```

## Next Steps

See the full [CLI reference](/docs/) for all commands, flags, and advanced usage.
