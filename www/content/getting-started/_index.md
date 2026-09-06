---
title: Getting Started
type: docs
---

This guide walks you through installing goGenerateCFToken, configuring your credentials, and generating your first Cloudflare API token.

## Prerequisites

- A [Cloudflare](https://www.cloudflare.com/) account with at least one zone
- A Cloudflare API token with the necessary permissions
- [Go](https://go.dev/) 1.27+ (if installing from source)

## Installation

### Install script

```bash
tmp=$(mktemp)
curl -sSfL https://raw.githubusercontent.com/nicholas-fedor/goGenerateCFToken/main/scripts/install.sh -o "$tmp" && sh "$tmp"
rm -f "$tmp"
```

On Linux this prefers a native `.deb` / `.rpm` / `.apk` / Arch package when sudo is available, otherwise it installs the release archive into `$HOME/go/bin`.

### From Source

```bash
go install github.com/nicholas-fedor/gogeneratecftoken@latest
```

### Pre-built Binary

Download the latest archive or distro package from the [GitHub releases page](https://github.com/nicholas-fedor/gogeneratecftoken/releases).

## Initial Setup

### 1. Initialize Configuration

Initialize the configuration file with your Cloudflare zone:

```bash
gogeneratecftoken config init
```

This prompts for the **Cloudflare zone name** (e.g., `example.com`) and writes
`$XDG_CONFIG_HOME/gogeneratecftoken/config.yaml` (typically `~/.config/gogeneratecftoken/config.yaml`).

Then store your Cloudflare API key. Interactive input is not echoed. On systems
without an OS keyring, the key is stored in a 0600 file next to the config:

```bash
gogeneratecftoken credentials set
```

Non-interactive alternatives:

```bash
export CF_API_TOKEN="your-master-api-token"
gogeneratecftoken credentials set --from-env

# Docker secrets / file-based credentials
gogeneratecftoken credentials set --from-file /run/secrets/cf_api_token
```

### 2. Verify Credentials

Test that your credentials are valid:

```bash
gogeneratecftoken credentials validate
```

## Configuration

Configuration is loaded from an XDG-compliant path:

```bash
~/.config/gogeneratecftoken/config.yaml
```

Override the path with `--config`:

```bash
gogeneratecftoken --config ./custom-config.yaml token generate myapp
```

## Basic Usage

### Managing Credentials

#### Store a New API Key

```bash
gogeneratecftoken credentials set
```

#### Remove Stored Credentials

```bash
gogeneratecftoken credentials remove
```

### Generating a Token

Generate a Cloudflare API token for a service:

```bash
gogeneratecftoken token generate myapp
```

This creates a token named `myapp.example.com` with:

- Zone read permissions
- DNS write permissions

The token is printed to stdout by default. Use `--output` to write it to a file:

```bash
gogeneratecftoken token generate myapp --output ./token.txt
```

#### Custom Token Name

Override the default naming convention:

```bash
gogeneratecftoken token generate myapp --name ci-deploy
```

#### Token Expiration

Set an expiration date for the token:

```bash
gogeneratecftoken token generate myapp --expires-on 2027-01-01T00:00:00Z
```

#### JSON Output

Output the token in JSON format for use in scripts:

```bash
gogeneratecftoken token generate myapp --json
```

### Managing Tokens

#### List Tokens

View all tokens associated with your Cloudflare account:

```bash
gogeneratecftoken token list
```

#### Revoke a Token

Revoke a token by its ID:

```bash
gogeneratecftoken token revoke abc123 --yes
```

## Updating

Update using the installation script:

```bash
tmp=$(mktemp)
curl -sSfL https://raw.githubusercontent.com/nicholas-fedor/goGenerateCFToken/main/scripts/install.sh -o "$tmp" && sh "$tmp" update
rm -f "$tmp"
```

## Uninstalling

Uninstall using the installation script:

```bash
tmp=$(mktemp)
curl -sSfL https://raw.githubusercontent.com/nicholas-fedor/goGenerateCFToken/main/scripts/install.sh -o "$tmp" && sh "$tmp" uninstall
rm -f "$tmp"
```

## Next Steps

See the full [CLI reference](/docs/) for all commands, flags, and advanced usage.
