<!-- markdownlint-disable -->
<div align="center">

# Cloudflare API Token Generator

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nicholas-fedor/gogeneratecftoken/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nicholas-fedor/gogeneratecftoken/tree/main)
[![codecov](https://codecov.io/gh/nicholas-fedor/gogeneratecftoken/branch/main/graph/badge.svg)](https://codecov.io/gh/nicholas-fedor/gogeneratecftoken)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1c48cfb7646d4009aa8c6f71287670b8)](https://www.codacy.com/gh/nicholas-fedor/gogeneratecftoken/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nicholas-fedor/gogeneratecftoken&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/nicholas-fedor/gogeneratecftoken?status.svg)](https://godoc.org/github.com/nicholas-fedor/gogeneratecftoken)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nicholas-fedor/go-remove)
[![Latest Version](https://img.shields.io/github/tag/nicholas-fedor/gogeneratecftoken.svg)](https://github.com/nicholas-fedor/gogeneratecftoken/releases)
[![Pulls from DockerHub](https://img.shields.io/docker/pulls/nickfedor/gogeneratecftoken.svg)](https://hub.docker.com/r/nickfedor/gogeneratecftoken)
[![AGPLv3 License](https://img.shields.io/github/license/nicholas-fedor/gogeneratecftoken.svg)](https://www.gnu.org/licenses/agpl-3.0)

----------

A simple CLI tool for generating and managing scoped Cloudflare API tokens.

</div>
<!-- markdownlint-restore -->

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
  - [Release Binaries](#release-binaries)
  - [Docker](#docker)
  - [Source](#source)
- [Usage](#usage)
  - [Overview](#overview)
  - [Configuration](#configuration)
    - [Configuration File](#configuration-file)
    - [Environment Variables](#environment-variables)
    - [CLI Flags](#cli-flags)
- [Contributing](#contributing)

## Quick Start

1. [Install](#installation) gogeneratecftoken

2. Create a Master API Key

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

3. Initialize the configuration:

    ```bash
    gogeneratecftoken config init
    ```

4. Add your zone when prompted:

    ```bash
    Enter Cloudflare zone name: <example.com>
    ```

5. Use `credentials set` to store the API key that you created in Step 1:

    ```bash
    Enter Cloudflare API key: <your API key>
    ```

6. Generate a token:

    ```bash
    gogeneratecftoken token generate test
    ```

    **Expected Output:**

    ```bash
    yoursuperlongandsecretserviceapitoken
    ```

## Installation

### Release Binaries

Download and install the latest binary for your platform from the [releases page](https://github.com/nicholas-fedor/gogeneratecftoken/releases).

The following are CLI scripts for installing to the user's `go/bin` directory:

- Windows (amd64):

    ```powershell
    New-Item -ItemType Directory -Path $HOME\go\bin -Force | Out-Null; iwr (iwr https://api.github.com/repos/nicholas-fedor/gogeneratecftoken/releases/latest | ConvertFrom-Json).assets.where({$_.name -like "*windows_amd64*.zip"}).browser_download_url -OutFile gogeneratecftoken.zip; Add-Type -AssemblyName System.IO.Compression.FileSystem; ($z=[System.IO.Compression.ZipFile]::OpenRead("$PWD\gogeneratecftoken.zip")).Entries | ? {$_.Name -eq 'gogeneratecftoken.exe'} | % {[System.IO.Compression.ZipFileExtensions]::ExtractToFile($_, "$HOME\go\bin\$($_.Name)", $true)}; $z.Dispose(); rm gogeneratecftoken.zip; if (Test-Path "$HOME\go\bin\gogeneratecftoken.exe") { Write-Host "Successfully installed gogeneratecftoken.exe to $HOME\go\bin" } else { Write-Host "Failed to install gogeneratecftoken.exe" }
    ```

- Linux (amd64):

    ```bash
    mkdir -p $HOME/go/bin && curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/gogeneratecftoken/releases/latest | grep -o 'https://[^"]*linux_amd64[^"]*\.tar\.gz') | tar -xz --strip-components=1 -C $HOME/go/bin gogeneratecftoken
    ```

- macOS (amd64):

    ```bash
    mkdir -p $HOME/go/bin && curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/gogeneratecftoken/releases/latest | grep -o 'https://[^"]*darwin_amd64[^"]*\.tar\.gz') | tar -xz --strip-components=1 -C $HOME/go/bin gogeneratecftoken
    ```

### Docker

Run gogeneratecftoken using Docker without installing it locally:

```bash
docker run --rm ghcr.io/nicholas-fedor/gogeneratecftoken:latest token generate test
```

> [!Warning]
> The Docker image runs as the `nobody` user (UID 65532). When mounting a configuration file, ensure proper file permissions:
>
> - The config file must be readable by UID 65532
> - Consider using `chmod 644` on the config file or mounting with appropriate user mapping
> - Environment variable credentials (`-e`) bypass file permission issues

To use a configuration file, mount it as a volume:

```bash
docker run --rm -v $HOME/.gogeneratecftoken/config.yaml:/config.yaml \
  ghcr.io/nicholas-fedor/gogeneratecftoken:latest token generate test --config /config.yaml
```

Alternatively, you can run as the current user to avoid permission issues:

```bash
docker run --rm -u $(id -u):$(id -g) -v $HOME/.gogeneratecftoken/config.yaml:/config.yaml \
  ghcr.io/nicholas-fedor/gogeneratecftoken:latest token generate test --config /config.yaml
```

Or pass credentials via environment variables:

```bash
docker run --rm -e CF_API_TOKEN="your-master-api-token" -e CF_ZONE="example.com" \
  ghcr.io/nicholas-fedor/gogeneratecftoken:latest token generate test
```

### Source

```bash
go install github.com/nicholas-fedor/gogeneratecftoken@latest
```

## Usage

### Overview

Invoke the program and use the `token generate` command with a subdomain as the argument to generate a Cloudflare API token.

The token will be named using the `subdomain.domain.tld` convention.

```bash
gogeneratecftoken token generate [SUBDOMAIN] [FLAGS]
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

`gogeneratecftoken` uses [Cobra](https://github.com/spf13/cobra). Configuration lives in an XDG-compliant YAML file. API keys are stored in the OS keyring (or a 0600 credential file when no keyring is available), not in the YAML.

#### Configuration File

Default location: `$XDG_CONFIG_HOME/gogeneratecftoken/config.yaml` (typically `~/.config/gogeneratecftoken/config.yaml`).

Legacy paths still searched: `~/.gogeneratecftoken/config.yaml` and `~/.goGenerateCFToken/config.yaml`.

Example:

```yaml
# gogeneratecftoken Configuration File
# https://github.com/nicholas-fedor/gogeneratecftoken

zone: "example.com"
account_id: ""      # optional
token_name: ""      # optional default token name
```

> [!Note]
> If using a custom configuration file location, then specify the location using the `--config` flag.
> Example:
>
> ```bash
> gogeneratecftoken token generate [SUBDOMAIN] --config [PATH]
> ```

#### Environment Variables

```bash
export CF_API_TOKEN="your-master-api-token"
export CF_API_TOKEN_FILE="/run/secrets/cf_api_token"
export CF_ZONE="example.com"
```

#### CLI Flags

You can use CLI flags directly instead of using a configuration file or setting environment variables.

- `t, --token`: Specify a master API token that has the permissions for creating additional tokens.
- `-z, --zone` : Specify a specific zone, i.e. example.com

## Contributing

Contributions to this project are welcomed.
Please see the [contributing documentation](/CONTRIBUTING.md) for more information.
