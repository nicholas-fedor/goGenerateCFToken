// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/prompt"
)

// newSetCommand creates the config set subcommand.
//
// Returns:
//   - *cobra.Command: The set command that updates the zone value in config.
func newSetCommand() *cobra.Command {
	cflags := &flags.ConfigFlags{}

	cmd := &cobra.Command{
		Use:   "set [zone]",
		Short: "Set configuration value",
		Long: `Set configuration values.

Zone may be given as an argument or entered interactively.
Use --account-id and --token-name to store generate defaults.`,
		Example: `  # Set zone with argument
  goGenerateCFToken config set example.com

  # Set account ID without changing zone
  goGenerateCFToken config set --account-id acc123

  # Interactive prompt
  goGenerateCFToken config set`,
		Args:    cobra.MaximumNArgs(1),
		GroupID: configGroup.ID,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetCmd(cmd, args, cflags)
		},
	}

	cflags.Bind(cmd.Flags())

	return cmd
}

// runSetCmd executes the config set command, prompting for zone if not provided.
//
// Parameters:
//   - cmd: Cobra command used for flag access.
//   - args: Command arguments, optionally containing the zone name.
//   - cflags: Config flags providing optional account ID and token name.
//
// Returns:
//   - error: Non-nil if zone cannot be read, is invalid, or config cannot be written.
func runSetCmd(cmd *cobra.Command, args []string, cflags *flags.ConfigFlags) error {
	cfgPath, err := resolveConfigPath(cmd)
	if err != nil {
		return err
	}

	err = config.EnsureFile(cfgPath)
	if err != nil {
		return fmt.Errorf("ensure config file: %w", err)
	}

	cfg, err := loadOrDefault(cfgPath)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		cfg.Zone = args[0]
	} else if cflags.AccountID == "" && cflags.TokenName == "" {
		zone, perr := prompt.Prompt("Enter Cloudflare zone name: ")
		if perr != nil {
			return fmt.Errorf("read zone: %w", perr)
		}

		cfg.Zone = zone
	}

	if cflags.AccountID != "" {
		cfg.AccountID = cflags.AccountID
	}

	if cflags.TokenName != "" {
		cfg.TokenName = cflags.TokenName
	}

	if cfg.Zone != "" && !config.ZonePattern.MatchString(cfg.Zone) {
		return fmt.Errorf("%w: %s", config.ErrInvalidZone, cfg.Zone)
	}

	log.Debug().Str("path", cfgPath).Str("zone", cfg.Zone).Msg("updating config")

	err = config.Write(cfgPath, *cfg)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	if cfg.Zone != "" {
		logging.Printf("Zone set: %s", cfg.Zone)
	}

	if cfg.AccountID != "" {
		logging.Printf("Account ID set: %s", cfg.AccountID)
	}

	if cfg.TokenName != "" {
		logging.Printf("Token name set: %s", cfg.TokenName)
	}

	return nil
}

// resolveConfigPath returns the config file path from the --config flag or
// Locate.
//
// Parameters:
//   - cmd: Cobra command used to read the config flag.
//
// Returns:
//   - string: The resolved configuration file path.
//   - error: Non-nil if the flag cannot be read or Locate fails.
func resolveConfigPath(cmd *cobra.Command) (string, error) {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return "", fmt.Errorf("get config flag: %w", err)
	}

	if cfgPath != "" {
		return cfgPath, nil
	}

	cfgPath, err = config.Locate()
	if err != nil {
		return "", fmt.Errorf("locate config: %w", err)
	}

	return cfgPath, nil
}

// loadOrDefault loads configuration from path, or returns defaults if zone is
// unset.
//
// Parameters:
//   - path: The configuration file path to load.
//
// Returns:
//   - *config.Config: The loaded or default configuration.
//   - error: Non-nil if the file cannot be loaded for reasons other than a
//     missing zone.
func loadOrDefault(path string) (*config.Config, error) {
	cfg, err := config.Load(path)
	if err == nil {
		return cfg, nil
	}

	if errors.Is(err, config.ErrZoneRequired) {
		d := config.Default()

		return &d, nil
	}

	return nil, fmt.Errorf("load existing config: %w", err)
}
