// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/logging"
	"github.com/nicholas-fedor/gogeneratecftoken/v2/internal/prompt"
)

// newInitCommand creates the config init subcommand.
//
// Returns:
//   - *cobra.Command: The init command configured for config setup.
func newInitCommand() *cobra.Command {
	cflags := &flags.ConfigFlags{}

	return &cobra.Command{
		Use:   "init [zone]",
		Short: "Initialize the config file",
		Long: `Initialize the configuration file.

Optionally accepts a zone as an argument. If not provided, prompts for it.
`,
		Example: `  # Interactive initialization
  goGenerateCFToken config init

  # Initialize with zone argument
  goGenerateCFToken config init example.com`,
		Args:    cobra.MaximumNArgs(1),
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInitCmd(cmd, args, cflags)
		},
	}
}

// runInitCmd executes the config init command, prompting for zone if not provided.
//
// Parameters:
//   - cmd: Cobra command used for flag access.
//   - args: Command arguments, optionally containing the zone name.
//   - _ cflags: Config flags (unused, present for consistency with other config commands).
//
// Returns:
//   - error: Non-nil if user input cannot be read, the zone is invalid,
//     or the config file cannot be written.
func runInitCmd(cmd *cobra.Command, args []string, _ *flags.ConfigFlags) error {
	var zone string
	if len(args) > 0 {
		zone = args[0]
	} else {
		z, err := prompt.Prompt("Enter Cloudflare zone name: ")
		if err != nil {
			return fmt.Errorf("read zone: %w", err)
		}

		zone = z
	}

	if !config.ZonePattern.MatchString(zone) {
		return fmt.Errorf("%w: %s", config.ErrInvalidZone, zone)
	}

	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("get config flag: %w", err)
	}

	if cfgPath == "" {
		cfgPath, err = config.Locate()
		if err != nil {
			return fmt.Errorf("locate config: %w", err)
		}
	}

	log.Debug().Str("path", cfgPath).Str("zone", zone).Msg("initializing config file")

	err = config.EnsureFile(cfgPath)
	if err != nil {
		return fmt.Errorf("ensure config file: %w", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Warn().Err(err).Msg("existing config is invalid - using defaults")

		defaultCfg := config.Default()
		cfg = &defaultCfg
	}

	cfg.Zone = zone

	err = config.Write(cfgPath, *cfg)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	logging.Printf("Initialization complete: %s (zone: %s). Use 'credentials set' to store the API key.", cfgPath, zone)

	return nil
}
