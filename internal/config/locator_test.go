// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"testing"
)

func Test_checkFilePermissions(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "file with secure permissions",
			path: "/tmp/secure.yaml",
		},
		{
			name: "file with insecure permissions",
			path: "/tmp/insecure.yaml",
		},
		{
			name: "nonexistent file",
			path: "/tmp/nonexistent.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkFilePermissions(tt.path)
		})
	}
}
