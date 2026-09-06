---
title: CLI Reference
type: docs
---

Complete command reference for goGenerateCFToken, organized by functional area.

## Configuration

View, set, reset, delete, and validate the configuration file.

| Command | Description |
|---------|-------------|
| [delete](/docs/config/delete/) | Delete the config file and directory |
| [init](/docs/config/init/) | Initialize the config file |
| [reset](/docs/config/reset/) | Reset the config file to defaults |
| [set](/docs/config/set/) | Set configuration value |
| [show](/docs/config/show/) | Show current configuration |
| [validate](/docs/config/validate/) | Validate configuration file |


## Credentials

Store, remove, and validate Cloudflare API credentials.

| Command | Description |
|---------|-------------|
| [remove](/docs/credentials/remove/) | Remove API key from OS keyring |
| [set](/docs/credentials/set/) | Store API key in the OS keyring or a local credential file |
| [validate](/docs/credentials/validate/) | Validate Cloudflare credentials |


## Tokens

Generate, list, get, and revoke Cloudflare API tokens.

| Command | Description |
|---------|-------------|
| [generate](/docs/token/generate/) | Generate a Cloudflare API token |
| [get](/docs/token/get/) | Get a single token's metadata |
| [list](/docs/token/list/) | List Cloudflare API tokens |
| [revoke](/docs/token/revoke/) | Revoke a Cloudflare API token |


## Version

Print the application version, including commit SHA and build details.

| Command | Description |
|---------|-------------|
| [version](/docs/version/) | Print the application version |


