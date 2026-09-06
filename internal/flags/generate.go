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
	// AccountID is the optional Cloudflare account ID for zone disambiguation.
	AccountID string
	// Output is the file path to write the token to instead of stdout.
	Output string
	// Format controls stdout output: text, json, or none.
	Format string
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
	// TTL is a human-friendly token lifetime (e.g. 24h, 7d). Converted to ExpiresOn.
	TTL string
}

// Bind attaches generate-specific flags to the provided flag set.
//
// Parameters:
//   - flags: The flag set to attach generate flags to.
func (gf *GenerateFlags) Bind(flags *pflag.FlagSet) {
	gf.bindAuth(flags)
	flags.StringVarP(
		&gf.Zone,
		"zone",
		"z",
		"",
		"Cloudflare zone name",
	)
	flags.StringVarP(
		&gf.Output,
		"output",
		"o",
		"",
		"Write token to file instead of stdout",
	)
	gf.bindJSON(flags)
	flags.StringVar(
		&gf.Name,
		"name",
		"",
		"Custom token name (default: service.zone)",
	)
	flags.BoolVar(
		&gf.DryRun,
		"dry-run",
		false,
		"Validate inputs without creating a token",
	)
	gf.bindTimeout(flags)
	flags.StringVar(
		&gf.ExpiresOn,
		"expires-on",
		"",
		"Token expiration date (RFC3339 format, e.g. 2027-01-01T00:00:00Z)",
	)
	flags.StringVar(
		&gf.TTL,
		"ttl",
		"",
		"Token lifetime as a human duration (e.g. 24h, 7d)",
	)
	flags.StringVar(
		&gf.AccountID,
		"account-id",
		"",
		"Cloudflare account ID for disambiguating zones across accounts",
	)
	flags.StringVar(
		&gf.Format,
		"format",
		"text",
		"Output format (text, json, none)",
	)
}

// BindGet attaches the subset of generate flags used by token get.
//
// Parameters:
//   - flags: The flag set to attach get flags to.
func (gf *GenerateFlags) BindGet(flags *pflag.FlagSet) {
	gf.bindAuth(flags)
	gf.bindJSON(flags)
	gf.bindTimeout(flags)
}

// bindAuth attaches the API token flag to the provided flag set.
//
// Parameters:
//   - flags: The flag set to attach the token flag to.
func (gf *GenerateFlags) bindAuth(flags *pflag.FlagSet) {
	flags.StringVarP(
		&gf.Token,
		"token",
		"t",
		"",
		"Cloudflare API token (prefer CF_API_TOKEN env var or keyring)",
	)
}

// bindJSON attaches the JSON output flag to the provided flag set.
//
// Parameters:
//   - flags: The flag set to attach the JSON flag to.
func (gf *GenerateFlags) bindJSON(flags *pflag.FlagSet) {
	flags.BoolVar(
		&gf.JSON,
		"json",
		false,
		"Output token in JSON format",
	)
}

// bindTimeout attaches the API timeout flag to the provided flag set.
//
// Parameters:
//   - flags: The flag set to attach the timeout flag to.
func (gf *GenerateFlags) bindTimeout(flags *pflag.FlagSet) {
	flags.IntVar(
		&gf.Timeout,
		"timeout",
		DefaultTimeout,
		"Timeout in seconds for API calls",
	)
}
