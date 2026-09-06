---
title: goGenerateCFToken
description: A CLI for Cloudflare API token management
---

## Introduction

This is a simple CLI tool for generating Cloudflare API tokens for use by tools, such as [Traefik](https://traefik.io/traefik), [Caddy](https://caddyserver.com/), or [Certbot](https://certbot.eff.org/).

It is intentionally limited to only creating tokens with `Zone:Zone:Read` and `Zone:DNS:Edit` permissions while also supporting basic token lifecycle management.

## Features

- **Token Generation** — Create API tokens with zone read and DNS write permissions
- **Token Management** — List and revoke existing tokens
- **Credential Storage** — Securely store API keys in the OS keyring
- **Configuration** — XDG-compliant config file support

## Quick Start

1. Install the `gogeneratecftoken` binary:

    ```bash
    go install github.com/nicholas-fedor/gogeneratecftoken@latest
    ```

2. Initialize the configuration and follow the prompt:

    ```bash
    gogeneratecftoken config init
    ```

3. Add your master API token:

    ```bash
    gogeneratecftoken credentials set
    ```

4. Validate your credentials:

    ```bash
    gogeneratecftoken credentials validate
    ```

5. Generate a token:

    ```bash
    gogeneratecftoken token generate myservice
    ```

6. List your tokens:

    ```bash
    gogeneratecftoken token list
    ```

## Documentation

See the [Documentation](/docs/) section for the full CLI reference.
