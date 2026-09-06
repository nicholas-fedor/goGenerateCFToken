---
title: CLI Reference
type: docs
---

Complete command reference for goGenerateCFToken, organized by functional area.

## Configuration

View, set, reset, delete, and validate the configuration file.

| Command | Description |
|---------|-------------|
| [delete](/cli-reference/config/delete/) | Delete the config file and directory |
| [init](/cli-reference/config/init/) | Initialize the config file |
| [reset](/cli-reference/config/reset/) | Reset the config file to defaults |
| [set](/cli-reference/config/set/) | Set configuration value |
| [show](/cli-reference/config/show/) | Show current configuration |
| [validate](/cli-reference/config/validate/) | Validate configuration file |


## Credentials

Store, remove, and validate Cloudflare API credentials.

| Command | Description |
|---------|-------------|
| [remove](/cli-reference/credentials/remove/) | Remove API key from OS keyring |
| [set](/cli-reference/credentials/set/) | Store API key in the OS keyring or a local credential file |
| [validate](/cli-reference/credentials/validate/) | Validate Cloudflare credentials |


## Tokens

Generate, list, get, and revoke Cloudflare API tokens.

| Command | Description |
|---------|-------------|
| [generate](/cli-reference/token/generate/) | Generate a Cloudflare API token |
| [get](/cli-reference/token/get/) | Get a single token's metadata |
| [list](/cli-reference/token/list/) | List Cloudflare API tokens |
| [revoke](/cli-reference/token/revoke/) | Revoke a Cloudflare API token |


## Version

Print the application version, including commit SHA and build details.

| Command | Description |
|---------|-------------|
| [version](/cli-reference/version/) | Print the application version |


