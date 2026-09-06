// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Error().Err(err).Msg("execution failed")
		os.Exit(1)
	}
}
