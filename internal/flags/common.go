// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"github.com/spf13/pflag"
)

const (
	// DefaultTimeout is the default API call timeout in seconds.
	DefaultTimeout = 30
)

// CommonFlags contains flags common to multiple commands.
type CommonFlags struct {
	// LogLevel sets the logging verbosity level (debug, info, warn, error).
	LogLevel string
	// Config is the path to the configuration file.
	Config string
	// Quiet suppresses all output except errors.
	Quiet bool
	// Verbose enables verbose output.
	Verbose bool
}

// Bind attaches the common flags to the provided flag set.
// The flags' values are stored in the CommonFlags receiver.
//
// Parameters:
//   - flags: The flag set to attach common flags to.
func (cf *CommonFlags) Bind(flags *pflag.FlagSet) {
	// --log-level, -l: Set the logging verbosity level
	flags.StringVarP(
		&cf.LogLevel,
		"log-level",
		"l",
		"info",
		"Set logging level (debug, info, warn, error)",
	)
	// --config, -c: Path to the configuration file
	flags.StringVarP(
		&cf.Config,
		"config",
		"c",
		"",
		"Path to config file",
	)
	// --quiet, -q: Suppress non-error output
	flags.BoolVarP(
		&cf.Quiet,
		"quiet",
		"q",
		false,
		"Suppress all output except errors",
	)
	// --verbose, -v: Enable detailed output
	flags.BoolVarP(
		&cf.Verbose,
		"verbose",
		"v",
		false,
		"Enable verbose output",
	)
}
