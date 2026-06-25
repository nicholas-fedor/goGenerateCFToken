// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"github.com/spf13/pflag"
)

// ConfigFlags contains flags for the config command.
type ConfigFlags struct {
	CommonFlags

	// Format specifies the output format for config show (yaml, json).
	Format string
}

// Bind attaches config-specific flags to the provided flag set.
//
// Parameters:
//   - flags: The flag set to attach config flags to.
func (cf *ConfigFlags) Bind(flags *pflag.FlagSet) {
}

// BindShow attaches flags for the config show subcommand.
//
// Parameters:
//   - flags: The flag set to attach the format flag to.
func (cf *ConfigFlags) BindShow(flags *pflag.FlagSet) {
	// --format, -f: Output format (yaml, json)
	flags.StringVarP(
		&cf.Format,
		"format",
		"f",
		"yaml",
		"Output format (yaml, json)",
	)
}
