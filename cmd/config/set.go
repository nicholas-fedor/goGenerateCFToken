// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
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
		Long: `Set a configuration value.

If zone is not provided, prompts for it. Currently only the zone field is supported.`,
		Example: `  # Set zone with argument
  goGenerateCFToken config set example.com

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
//   - _ cflags: Config flags (unused, present for consistency).
//
// Returns:
//   - error: Non-nil if zone cannot be read, is invalid, or config cannot be written.
func runSetCmd(cmd *cobra.Command, args []string, _ *flags.ConfigFlags) error {
	var zone string
	if len(args) > 0 {
		zone = args[0]
	} else {
		z, err := promptZone("Enter Cloudflare zone name: ")
		if err != nil {
			return fmt.Errorf("read zone: %w", err)
		}

		zone = z
	}

	if !config.ZonePattern.MatchString(zone) {
		return fmt.Errorf("invalid zone name: %s", zone)
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

	log.Debug().Str("path", cfgPath).Str("zone", zone).Msg("updating config zone")

	cfg := config.Config{Zone: zone}

	err = config.EnsureFile(cfgPath)
	if err != nil {
		return fmt.Errorf("ensure config file: %w", err)
	}

	err = config.Write(cfgPath, cfg)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	logging.Printf("Zone set: %s", zone)

	return nil
}

// promptZone displays a message to stderr and reads a zone name from stdin.
//
// Parameters:
//   - message: The prompt text to display.
//
// Returns:
//   - string: The zone name input (trimmed of trailing newline).
//   - error: Non-nil if reading from stdin fails or input is empty.
func promptZone(message string) (string, error) {
	fmt.Fprint(os.Stderr, message)

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	result := strings.TrimSpace(input)
	if result == "" {
		return "", errors.New("zone cannot be empty")
	}

	return result, nil
}