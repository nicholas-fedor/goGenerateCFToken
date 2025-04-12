// Package config manages the configuration for the goGenerateCFToken CLI tool.
//
// The goGenerateCFToken tool relies on a YAML configuration file, environment
// variables, or command-line flags to specify the Cloudflare API token and zone
// name. This package uses Viper to handle configuration loading from a file
// (default: ~/.goGenerateCFToken/config.yaml), environment variables (prefixed
// with CF_, e.g., CF_API_TOKEN), or defaults.
//
// The configuration process involves:
//  1. Setting the configuration file path, either explicitly via a flag or
//     defaulting to the user’s home directory.
//  2. Binding environment variables with a "CF" prefix and replacing dots with
//     underscores (e.g., api_token becomes CF_API_TOKEN).
//  3. Loading the configuration file, reporting errors or success to stderr.
//
// The package defines a Viper interface to abstract configuration operations,
// allowing for dependency injection during testing. Key functions include
// InitConfig to start the configuration process, and internal helpers to set
// file paths, environment variables, and load the config.
package config
