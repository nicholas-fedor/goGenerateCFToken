// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package config manages configuration for the goGenerateCFToken CLI tool.
//
// The configuration is loaded from an XDG-compliant config file
// ($XDG_CONFIG_HOME/gogeneratecftoken/config.yaml) or the current directory.
// Configuration file permissions are checked on load and a warning is logged if
// the file is accessible by group or others (permissions > 0600).
package config
