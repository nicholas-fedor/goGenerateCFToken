// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package flags

import (
	"github.com/spf13/pflag"
)

// VersionFlags contains flags for the version command.
type VersionFlags struct {
	CommonFlags

	// JSON outputs version information in JSON format.
	JSON bool
}

// Bind attaches version-specific flags to the provided flag set.
// Common flags (log-level, config, quiet, verbose) are inherited from the root
// command's persistent flags and are not re-bound here.
//
// Parameters:
//   - flags: The flag set to attach the JSON flag to.
func (vf *VersionFlags) Bind(flags *pflag.FlagSet) {
	// --json: Output version information in JSON format
	flags.BoolVar(
		&vf.JSON,
		"json",
		false,
		"Output version information in JSON format",
	)
}
