package version

import (
	"runtime/debug"
	"time"
)

// Constants for repeated string values.
const (
	devVersion   = "dev"
	unknownValue = "unknown"
)

// These values are populated by GoReleaser during release builds.
var (
	// Version is the CLI tool's version (e.g., "v0.0.1").
	Version = devVersion
	// Commit is the Git commit SHA (e.g., "abc123").
	Commit = unknownValue
	// Date is the build or commit timestamp in RFC3339 format (e.g., "2025-05-07T00:00:00Z").
	Date = unknownValue
)

// Info holds version information for the CLI.
type Info struct {
	Version string
	Commit  string
	Date    string
}

// GetVersionInfo returns version information, using debug.ReadBuildInfo for source builds
// or GoReleaser variables for release builds.
func GetVersionInfo() Info {
	version := Version
	commit := Commit
	date := Date

	// If building from source (not GoReleaser), try to get version info from debug.ReadBuildInfo
	if version == devVersion || version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			// Get the module version (e.g., v1.1.4 or v1.1.4+dirty)
			version = info.Main.Version
			if version == "(devel)" || version == "" {
				version = unknownValue
			}

			// Extract VCS information (Git commit and timestamp)
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					commit = setting.Value
				case "vcs.time":
					if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
						date = t.Format(time.RFC3339)
					}
				case "vcs.modified":
					if setting.Value == "true" && version != unknownValue &&
						!contains(version, "+dirty") {
						version += "+dirty"
					}
				}
			}
		}
	} else {
		// GoReleaser provides a valid version without 'v' prefix, so add it
		if version != "v" {
			version = "v" + version
		}
	}

	// Fallback defaults if still unset or invalid
	if version == "" || version == devVersion || version == "v" {
		version = unknownValue
	}
	if commit == "" {
		commit = unknownValue
	}
	if date == "" {
		date = unknownValue
	}

	return Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
