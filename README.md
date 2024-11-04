# Cloudflare API Token Generator

This is a simple tool written in Go for generating Cloudflare API tokens.
The generated tokens are intended to be service-specific, i.e. Plex, Radarr, etc, for generating SSL certificates by tools, such as Certbot.

## Prerequisites

Create a `.env` file in the directory where the program will be run.
Include the following content:

```console
CLOUDFLARE_API_TOKEN = "[Cloudflare API Token]"
ZONE = "[example.com]"
```

### Cloudflare API Token

An API token needs to be manually created to be used by the application itself via the [Cloudflare account dashboard](https://dash.cloudflare.com/profile/api-tokens).

1) Select the `Create Additional Tokens` template
2) Update the `Token name`, as needed
3) Add the following permission: `Zone` - `Zone` - `Read`
4) Zone Resources: `Include` - `Specific Zone` - `[example.com]`
5) Client IP Address Filtering: This limits usage of the token to a specified IP address
6) Add the newly created token to a `.env` file

### Zone

Specify the particular zone, i.e. `example.com`, that is linked to the API token.

## Usage

1) Ensure the prerequisite `.env` file is setup

2) Launch the tool via the console from within the same directory as the `.env` file and specify the service name as a command-line parameter:

    ```console
    goGenerateCFToken [service name]
    ```

3) The generated token will be printed to the console

## Resources

- [Cloudflare's API Endpoints](https://developers.cloudflare.com/api-next)
- [Cloudflare's Go SDK](https://github.com/cloudflare/cloudflare-go)
- [Joho's DotENV Library](https://github.com/joho/godotenv)
