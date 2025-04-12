// Package cmd provides the command-line interface for the goGenerateCFToken tool.
//
// The goGenerateCFToken CLI generates Cloudflare API tokens with DNS edit permissions.
// It supports a single command, "generate", which creates a token based on a provided
// service name and configuration settings (API token and zone name). The tool uses
// Cobra for command handling and Viper for configuration management.
//
// The root command, "goGenerateCFToken", initializes the CLI and supports a persistent
// --config flag to specify the configuration file path. The configuration file contains
// the Cloudflare API token and zone name, which can also be set via flags (--token, --zone)
// or environment variables (CF_API_TOKEN, CF_ZONE).
//
// Example usage:
//
//	goGenerateCFToken generate service
//
// This command generates a token named "service.example.com" for the zone specified in
// the configuration, printing the token value to stdout.
//
// The package defines error constants for common failure cases, such as missing configuration
// or client initialization errors, ensuring clear error reporting.
package cmd
