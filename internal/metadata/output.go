// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// PrintDefault writes the single-line version string to writer.
//
// Parameters:
//   - writer: The output writer to write the version string to.
//
// Returns:
//   - error: Non-nil if writing to the writer fails.
func PrintDefault(writer io.Writer) error {
	_, err := io.WriteString(writer, Name+" "+String()+"\n")
	if err != nil {
		return fmt.Errorf("write default output: %w", err)
	}

	return nil
}

// PrintVerbose writes multi-line detailed version info to writer.
//
// Parameters:
//   - writer: The output writer to write the version details to.
//
// Returns:
//   - error: Non-nil if writing to the writer fails.
func PrintVerbose(writer io.Writer) error {
	info := GetInfo()

	lines := []string{
		"name:      " + info.Name,
		"version:   " + info.Version,
	}
	if info.CommitSHA != emptyString {
		lines = append(lines, "commitSha: "+info.CommitSHA)
	}

	if info.BuildTime != emptyString {
		lines = append(lines, "buildTime: "+info.BuildTime)
	}

	lines = append(lines,
		"goVersion: "+info.GoVersion,
		"os:        "+info.OS,
		"arch:      "+info.Arch,
	)

	content := strings.Join(lines, "\n") + "\n"

	_, err := io.WriteString(writer, content)
	if err != nil {
		return fmt.Errorf("write verbose output: %w", err)
	}

	return nil
}

// PrintJSON writes version info as indented JSON to writer.
//
// Parameters:
//   - writer: The output writer to encode JSON to.
//
// Returns:
//   - error: Non-nil if JSON encoding fails.
func PrintJSON(writer io.Writer) error {
	info := GetInfo()

	enc := json.NewEncoder(writer)
	enc.SetIndent("", "  ")

	err := enc.Encode(info)
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
