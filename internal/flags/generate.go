// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"github.com/spf13/pflag"
)

// GenerateFlags contains flags for the generate command.
type GenerateFlags struct {
	CommonFlags

	// Token is the Cloudflare API token override. Prefer CF_API_TOKEN_FILE, which keeps
	// the secret out of the argument list; see bindAuth for why this flag is deprecated.
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
// The flag is marked deprecated because a secret passed as an argument is exposed to every
// local user through process listings, is readable from /proc/<pid>/cmdline by any process
// running as the same user, and is persisted in shell history and in any CI log that echoes
// the command. It keeps working, but pflag hides it from help once marked, so the flag is
// no longer advertised. A local file supplied through CF_API_TOKEN_FILE keeps the secret out
// of both the argument list and the environment.
//
// Parameters:
//   - flags: The flag set to attach the token flag to.
func (gf *GenerateFlags) bindAuth(flags *pflag.FlagSet) {
	flags.StringVarP(
		&gf.Token,
		"token",
		"t",
		"",
		"Cloudflare API token (prefer CF_API_TOKEN_FILE or the OS keyring)",
	)

	// The error is unreachable: the flag was registered on this same set immediately above,
	// and MarkDeprecated only fails when the name is absent. bindAuth has no error return and
	// giving it one would change the signatures of Bind and BindGet for an impossible case.
	_ = flags.MarkDeprecated(
		"token",
		"the --token flag exposes the API key through process listings, shell history and CI logs; "+
			"set CF_API_TOKEN_FILE, set CF_API_TOKEN, or run 'credentials set' instead",
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
