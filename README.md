<!-- markdownlint-disable -->
<div align="center">

# Cloudflare API Token Generator

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nicholas-fedor/goGenerateCFToken/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nicholas-fedor/goGenerateCFToken/tree/main)
[![codecov](https://codecov.io/gh/nicholas-fedor/goGenerateCFToken/branch/main/graph/badge.svg)](https://codecov.io/gh/nicholas-fedor/goGenerateCFToken)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1c48cfb7646d4009aa8c6f71287670b8)](https://www.codacy.com/gh/nicholas-fedor/goGenerateCFToken/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nicholas-fedor/goGenerateCFToken&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicholas-fedor/goGenerateCFToken)](https://goreportcard.com/report/github.com/nicholas-fedor/goGenerateCFToken)
[![GoDoc](https://godoc.org/github.com/nicholas-fedor/gogeneratecftoken?status.svg)](https://godoc.org/github.com/nicholas-fedor/gogeneratecftoken)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nicholas-fedor/go-remove)
[![latest version](https://img.shields.io/github/tag/nicholas-fedor/goGenerateCFToken.svg)](https://github.com/nicholas-fedor/goGenerateCFToken/releases)
[![AGPLv3 License](https://img.shields.io/github/license/nicholas-fedor/goGenerateCFToken.svg)](https://www.gnu.org/licenses/agpl-3.0)

----------

A simple CLI tool for generating Cloudflare API tokens for use by tools, such as [Traefik](https://traefik.io/traefik), [Caddy](https://caddyserver.com/), or [Certbot](https://certbot.eff.org/)

</div>
<!-- markdownlint-restore -->

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
  - [Release Binaries](#release-binaries)
  - [Source](#source)
- [Usage](#usage)
  - [Overview](#overview)
  - [Configuration](#configuration)
- [Contributing](#contributing)

## Quick Start

1. [Install](#installation) goGenerateCFToken

2. Create a Master API Token

    <!-- markdownlint-disable -->
    <ol type="a">
      <li>Go to your <a href="https://dash.cloudflare.com/profile/api-tokens">Cloudflare Dashboard</a></li>
      <li>Create a token with following permissions:</li>
            <ul>
                <li>Zone > Zone > Read</li>
                <li>User > API Tokens > Edit</li>
                <li>Include > Specific Zone > example.com</li>
            </ul>
      <li>Save/copy the token</li>
    </ol>
    <!-- markdownlint-restore -->

3. Setup the Configuration File

    - Download the configuration file template to `$HOME/.goGenerateCFToken/config.yaml`:

      - Windows

        ```powershell
        New-Item -ItemType Directory -Path $HOME\.goGenerateCFToken -Force; iwr -Uri https://github.com/nicholas-fedor/   goGenerateCFToken/raw/main/config.yaml.template -OutFile $HOME\.goGenerateCFToken\config.yaml
        ```

      - Linux

        ```bash
        mkdir -p $HOME/.goGenerateCFToken && curl -L https://github.com/nicholas-fedor/goGenerateCFToken/raw/main/config.   yaml.template -o $HOME/.goGenerateCFToken/config.yaml
        ```

      - macOS

        ```bash
        mkdir -p $HOME/.goGenerateCFToken && curl -L https://github.com/nicholas-fedor/goGenerateCFToken/raw/main/config.   yaml.template -o $HOME/.goGenerateCFToken/config.yaml
        ```

    - Edit `$HOME/.goGenerateCFToken/config.yaml` to add your master API token and zone:

    ```yaml
    api_token: "your-master-api-token"
    zone: "example.com"
    ```

4. Generate a Test Token

    ```bash
    goGenerateCFToken generate test
    ```

    **Expected Output:**

    ```bash
    Generating API token: test.example.com
    yoursuperlongandsecretserviceapitoken
    ```

## Installation

### Release Binaries

Download and install the latest binary for your platform from the [releases page](https://github.com/nicholas-fedor/goGenerateCFToken/releases).

The following are CLI scripts for installing to the user's `go/bin` directory:

- Windows (amd64):

    ```powershell
    New-Item -ItemType Directory -Path $HOME\go\bin -Force | Out-Null; iwr (iwr https://api.github.com/repos/nicholas-fedor/goGenerateCFToken/releases/latest | ConvertFrom-Json).assets.where({$_.name -like "*windows_amd64*.zip"}).browser_download_url -OutFile goGenerateCFToken.zip; Add-Type -AssemblyName System.IO.Compression.FileSystem; ($z=[System.IO.Compression.ZipFile]::OpenRead("$PWD\goGenerateCFToken.zip")).Entries | ? {$_.Name -eq 'goGenerateCFToken.exe'} | % {[System.IO.Compression.ZipFileExtensions]::ExtractToFile($_, "$HOME\go\bin\$($_.Name)", $true)}; $z.Dispose(); rm goGenerateCFToken.zip; if (Test-Path "$HOME\go\bin\goGenerateCFToken.exe") { Write-Host "Successfully installed goGenerateCFToken.exe to $HOME\go\bin" } else { Write-Host "Failed to install goGenerateCFToken.exe" }
    ```

- Linux (amd64):

    ```bash
    mkdir -p $HOME/go/bin && curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/goGenerateCFToken/releases/latest | grep -o 'https://[^"]*linux_amd64[^"]*\.tar\.gz') | tar -xz --strip-components=1 -C $HOME/go/bin goGenerateCFToken
    ```

- macOS (amd64):

    ```bash
    mkdir -p $HOME/go/bin && curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/goGenerateCFToken/releases/latest | grep -o 'https://[^"]*darwin_amd64[^"]*\.tar\.gz') | tar -xz --strip-components=1 -C $HOME/go/bin goGenerateCFToken
    ```

### Source

```bash
go install github.com/nicholas-fedor/gogeneratecftoken@latest
```

## Usage

### Overview

Invoke the program and use the `generate` command with a subdomain as the argument to generate a Cloudflare API token.

The token will be named using the `subdomain.domain.tld` convention.

```bash
goGenerateCFToken generate [SUBDOMAIN] [FLAGS]
```

| Flags         | Input Type | Description                               |
|---------------|------------|-------------------------------------------|
| `--config`    | String     | Specify a configuration file location     |
| `-t, --token` | String     | Specify a Cloudflare API master token     |
| `-z, --zone`  | String     | Specify a domain name, i.e. example.com   |
| `-h, --help`  | None       | Show the help information for the command |

> [!Warning]
> The Cloudflare API token will only be shown via the standard output. Remember to save it in a secure location!

### Configuration

In order to generate Cloudflare API tokens, the program requires the following:

- A master API token with the permissions to generate additional Cloudflare API tokens
- A zone (i.e. example.com)
- A sudomain (i.e. "test" from test.example.com)

`goGenerateCFToken` uses [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) to enable configuration functionality.

### Configuration File

Default Location: `$HOME/.goGenerateCFToken/config.yaml`

Example:

```yaml
# goGenerateCFToken Configuration File
# https://github.com/nicholas-fedor/goGenerateCFToken

# https://dash.cloudflare.com/profile/api-tokens
# Token Name: [Add your token name here for reference]
# Permissions: Zone: Read & API Tokens: Edit
api_token: "your-cloudflare-api-token-here"
zone: "example.com"
```

> [!Note]
> If using a custom configuration file location, then specify the location using the `--config` flag.
> Example:
>
> ```bash
> goGenerateCFToken [SUBDOMAIN] --config [PATH]
> ```

### Environment Variables

If no config file is found or specified, then the program falls back to environment variables.

```bash
export CF_API_TOKEN="your-master-api-token"
export CF_ZONE="example.com"
```

### CLI Flags

You can use CLI flags directly instead of using a configuration file or setting environment variables.

- `t, --token`: Specify a master API token that has the permissions for creating additional tokens.
- `-z, --zone` : Specify a specific zone, i.e. example.com

## Contributing

Contributions to this project are welcomed.
Please see the [contributing documentation](/CONTRIBUTING.md) for more information.
