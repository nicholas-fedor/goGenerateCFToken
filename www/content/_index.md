---
title: goGenerateCFToken
description: A CLI for Cloudflare API token management
---

A CLI for creating and managing Cloudflare API tokens with DNS edit permissions.

## Features

- **Token Generation** — Create API tokens with zone read and DNS write permissions
- **Token Management** — List and revoke existing tokens
- **Credential Storage** — Securely store API keys in the OS keyring
- **Configuration** — XDG-compliant config file support

## Quick Start

```bash
# Install
go install github.com/nicholas-fedor/gogeneratecftoken@latest

# Initialize configuration
goGenerateCFToken config init

# Generate a token
goGenerateCFToken token generate myservice

# List tokens
goGenerateCFToken token list

# Validate credentials
goGenerateCFToken credentials validate
```

## Documentation

See the [Documentation](/docs/) section for the full CLI reference.
