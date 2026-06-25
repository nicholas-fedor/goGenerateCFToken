// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"github.com/spf13/pflag"
)

// GenerateFlags contains flags for the generate command.
type GenerateFlags struct {
	CommonFlags

	// Token is the Cloudflare API token override (prefer CF_API_TOKEN env var or keyring).
	Token string
	// Zone is the Cloudflare zone name (overrides config file).
	Zone string
	// Output is the file path to write the token to instead of stdout.
	Output string
	// JSON enables JSON output format.
	JSON bool
	// Name overrides the default token name (service.zone).
	Name string
	// DryRun validates inputs without creating a token.
	DryRun bool
	// Timeout is the API call timeout in seconds.
	Timeout int
	// ExpiresOn is the token expiration date in RFC3339 format.
	ExpiresOn string
}

// Bind attaches generate-specific flags to the provided flag set.
//
// Parameters:
//   - flags: The flag set to attach generate flags to.
func (gf *GenerateFlags) Bind(flags *pflag.FlagSet) {
	// --token, -t: Cloudflare API token (prefer CF_API_TOKEN env var or keyring)
	flags.StringVarP(
		&gf.Token,
		"token",
		"t",
		"",
		"Cloudflare API token (prefer CF_API_TOKEN env var or keyring)",
	)
	// --zone, -z: Cloudflare zone name
	flags.StringVarP(
		&gf.Zone,
		"zone",
		"z",
		"",
		"Cloudflare zone name",
	)
	// --output, -o: Write token to file instead of stdout
	flags.StringVarP(
		&gf.Output,
		"output",
		"o",
		"",
		"Write token to file instead of stdout",
	)
	// --json: Output token in JSON format
	flags.BoolVar(
		&gf.JSON,
		"json",
		false,
		"Output token in JSON format",
	)
	// --name: Custom token name (default: service.zone)
	flags.StringVar(
		&gf.Name,
		"name",
		"",
		"Custom token name (default: service.zone)",
	)
	// --dry-run: Validate inputs without creating a token
	flags.BoolVar(
		&gf.DryRun,
		"dry-run",
		false,
		"Validate inputs without creating a token",
	)
	// --timeout: Timeout in seconds for API calls
	flags.IntVar(
		&gf.Timeout,
		"timeout",
		DefaultTimeout,
		"Timeout in seconds for API calls",
	)
	// --expires-on: Token expiration date (RFC3339 format, e.g. 2027-01-01T00:00:00Z)
	flags.StringVar(
		&gf.ExpiresOn,
		"expires-on",
		"",
		"Token expiration date (RFC3339 format, e.g. 2027-01-01T00:00:00Z)",
	)
}
