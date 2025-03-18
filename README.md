<div align="center">

# Cloudflare API Token Generator

A simple CLI tool for generating Cloudflare API tokens.

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nicholas-fedor/goGenerateCFToken/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nicholas-fedor/goGenerateCFToken/tree/main)
[![codecov](https://codecov.io/gh/nicholas-fedor/goGenerateCFToken/branch/main/graph/badge.svg)](https://codecov.io/gh/nicholas-fedor/goGenerateCFToken)
[![GoDoc](https://godoc.org/github.com/nicholas-fedor/goGenerateCFToken?status.svg)](https://godoc.org/github.com/nicholas-fedor/goGenerateCFToken)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicholas-fedor/goGenerateCFToken)](https://goreportcard.com/report/github.com/nicholas-fedor/goGenerateCFToken)
[![latest version](https://img.shields.io/github/tag/nicholas-fedor/goGenerateCFToken.svg)](https://github.com/nicholas-fedor/goGenerateCFToken/releases)
[![AGPLv3 License](https://img.shields.io/github/license/nicholas-fedor/goGenerateCFToken.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1c48cfb7646d4009aa8c6f71287670b8)](https://www.codacy.com/gh/nicholas-fedor/goGenerateCFToken/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nicholas-fedor/goGenerateCFToken&amp;utm_campaign=Badge_Grade)

----------

</div>

## Overview

The generated tokens are intended to be service-specific, i.e. Plex, Radarr, etc, for generating SSL certificates by tools, such as [Certbot](https://certbot.eff.org/) or [Lego](https://go-acme.github.io/lego/).

## Prerequisites

### Cloudflare-managed Domain Name

You will need to own a domain name configured to use Cloudflare's nameservers.  
Once that is setup, you will be able to specify the domain's zone (i.e. `example.com`) that you want linked to the API token.

### Cloudflare API Token with `API Tokens: Edit` Permissions

An API token needs to be manually created via the [Cloudflare account dashboard](https://dash.cloudflare.com/profile/api-tokens) to be used by `goGenerateCFToken`.

1) Select the `Create Additional Tokens` template
2) Update the `Token name`, as needed
3) Add the following permission: `Zone` - `Zone` - `Read`
4) Zone Resources: `Include` - `Specific Zone` - `example.com`
5) Client IP Address Filtering: This limits usage of the token to a specified IP address
6) Save the newly-created API token to a safe location

## Usage

For access to usage instructions:

```console
goGenerateCFToken -h
```

### Configuration

There are several options for providing the `api_token` and `zone` to `goGenerateCFToken`, as it uses [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) to enable configuration functionality.

- YAML Config File:  
    Setup a configuration file, as shown in the provided `config.yaml.template`, in a location, such as the default `$HOME/config.yaml`.  
    If using a custom location, then use the `--config [path]` flag.
- Environment Variables:  
    If no config file is found or specified, then the program falls back to environment variables.

    ```console
    export CF_API_TOKEN="your-cloudflare-api-token-here"
    export CF_ZONE="example.com"
    ```

- CLI Flags:

    To see the available flags:

    ```console
    goGenerateCFToken generate -h
    ```

    Example Usage:

    ```console
    goGenearteCFToken generate [service name] -z [example.com] -t [supersecretcftoken]
    ```

### Recommended Usage

1) Copy the template to `$HOME/config.yaml`:

    ```console
    cp ./config.yaml.template $HOME/config.yaml
    ```

2) Obtain your master Cloudflare API token for creating additional service-specific tokens

3) Save the master token and associated zone to the configuration file

4) Generate a test API token to confirm configuration is successful:

    ```console
    goGenerateCFToken generate test
    ```

5) If successful, you should see the following output:

    ```console
    Generating API token: test.example.com
    yoursuperlongandsecretserviceapitoken
    ```

## Resources

- [Cloudflare's API Endpoints](https://developers.cloudflare.com/api-next)
- [Cloudflare's Go SDK](https://github.com/cloudflare/cloudflare-go)
- [spf13/Cobra](https://github.com/spf13/cobra)
- [spf13/Viper](https://github.com/spf13/viper)

## Development

### GitHub releases

Using [GoReleaser](https://github.com/goreleaser/goreleaser-action) to build the release files.

[Quick Start](https://goreleaser.com/quick-start/)

To run a new build, update the tag, as follows:

```console
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```
