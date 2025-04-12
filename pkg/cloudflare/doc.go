// Package cloudflare provides functionality to interact with the Cloudflare API
// for generating API tokens with DNS edit permissions.
//
// The package defines a Client type that wraps the Cloudflare SDK client,
// implementing methods to list zones and create API tokens. It uses an APIInterface
// to abstract API calls, enabling dependency injection for testing. The main
// functionality includes retrieving zone IDs by name and generating tokens for
// specific services and zones, with permissions for zone read and DNS write access.
//
// Key components:
// - Client: Wraps the Cloudflare SDK client, providing methods for zone and token operations.
// - APIInterface: Defines methods for listing zones and creating tokens, used for mocking in tests.
// - GenerateToken: Creates a token with specified permissions for a given zone and service name.
// - GetZoneID: Retrieves a zone ID by name, handling cases for zero or multiple matches.
//
// The package includes error constants for common failure cases, such as missing
// credentials, uninitialized clients, or API errors, ensuring clear error reporting.
// Configuration values (e.g., API token) are expected to be provided via the config package.
package cloudflare
