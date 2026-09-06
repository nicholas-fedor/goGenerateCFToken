// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package version provides the CLI command for displaying application version information.
package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/metadata"
)

// NewCommand creates the version command.
//
// Returns:
//   - *cobra.Command: The version command that prints application version information.
func NewCommand() *cobra.Command {
	vflags := &flags.VersionFlags{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		Long:  "Print the application version, including commit SHA and build details.",
		Example: `  # Print version
  goGenerateCFToken version

  # Print detailed version info
  goGenerateCFToken version --verbose

  # Print version as JSON
  goGenerateCFToken version --json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionCmd(cmd, vflags)
		},
	}

	vflags.Bind(cmd.Flags())

	return cmd
}

// runVersionCmd executes the version command.
//
// Parameters:
//   - cmd: Cobra command used for output.
//   - vflags: Version flags controlling JSON output and verbose mode.
//
// Returns:
//   - error: Non-nil if the printer cannot be selected or output fails.
func runVersionCmd(cmd *cobra.Command, vflags *flags.VersionFlags) error {
	printer, err := selectPrinter(cmd, vflags)
	if err != nil {
		return fmt.Errorf("select printer: %w", err)
	}

	err = printer(cmd.OutOrStdout())
	if err != nil {
		return fmt.Errorf("print version: %w", err)
	}

	return nil
}

// selectPrinter returns the appropriate version output printer based on flags.
//
// Parameters:
//   - cmd: Cobra command providing access to the verbose flag.
//   - vflags: Version flags. If JSON is true, returns the JSON printer.
//
// Returns:
//   - func(io.Writer) error: The selected printer function.
//   - error: Non-nil if the verbose flag cannot be retrieved.
func selectPrinter(cmd *cobra.Command, vflags *flags.VersionFlags) (func(io.Writer) error, error) {
	if vflags.JSON {
		return metadata.PrintJSON, nil
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return nil, fmt.Errorf("get verbose flag: %w", err)
	}

	if verbose {
		return metadata.PrintVerbose, nil
	}

	return metadata.PrintDefault, nil
}
