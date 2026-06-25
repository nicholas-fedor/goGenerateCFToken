// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

// newGenerateCommand creates the token generate subcommand.
//
// Returns:
//   - *cobra.Command: The generate command that creates a Cloudflare API token
//     with DNS edit permissions for the specified service.
func newGenerateCommand() *cobra.Command {
	gflags := &flags.GenerateFlags{}

	cmd := &cobra.Command{
		Use:   "generate <service-name>",
		Short: "Generate a Cloudflare API token",
		Long:  "Generate a new Cloudflare API token with DNS edit permissions for the specified service.\n\nThe token name defaults to `service-name.zone`. Use --name to override.\nThe zone is loaded from the config file or --zone flag.\nThe API token is resolved from CF_API_TOKEN env var or the OS keyring.",
		Example: `  # Generate a token using the zone from config
  goGenerateCFToken token generate myapp

  # Generate with a custom zone
  goGenerateCFToken token generate myapp --zone example.com

  # Generate with a custom name and expiration
  goGenerateCFToken token generate myapp --name ci-deploy --expires-on 2027-01-01T00:00:00Z

  # Output as JSON
  goGenerateCFToken token generate myapp --json`,
		GroupID: tokenGroup.ID,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return loadConfig(cmd, gflags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerateCmd(cmd, args, gflags)
		},
	}

	gflags.Bind(cmd.Flags())

	return cmd
}

// loadConfig loads the configuration and sets the zone flag if not already provided.
//
// Parameters:
//   - cmd: Cobra command for retrieving the config flag value.
//   - gflags: Generate flags providing the optional zone; Zone may be updated.
//
// Returns:
//   - error: Non-nil if the configuration cannot be loaded.
func loadConfig(cmd *cobra.Command, gflags *flags.GenerateFlags) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("get config flag: %w", err)
	}

	log.Debug().Str("config_path", configPath).Msg("loading configuration")

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if gflags.Zone == "" {
		gflags.Zone = cfg.Zone
		log.Debug().Str("zone", gflags.Zone).Msg("using zone from config")
	}

	return nil
}

// resolveAPIKey returns the API key from the --token flag or the credential resolver.
//
// Parameters:
//   - tokenFlag: The value of the --token flag, if set.
//
// Returns:
//   - string: The resolved API key.
//   - error: Non-nil if no API key is available from any source.
func resolveAPIKey(tokenFlag string) (string, error) {
	if tokenFlag != "" {
		log.Debug().Msg("using API key from --token flag")

		return tokenFlag, nil
	}

	key, err := credentials.ResolveAPIKey()
	if err != nil {
		return "", fmt.Errorf("resolve API key: %w", err)
	}

	return key, nil
}

// runGenerateCmd executes the token generate command.
//
// Parameters:
//   - cmd: Cobra command providing context with timeout.
//   - args: Positional arguments; args[0] is the service name.
//   - gflags: Generate flags controlling zone, token, output format, and dry-run behavior.
//
// Returns:
//   - error: Non-nil if validation fails, the API key cannot be resolved,
//     or token generation fails.
func runGenerateCmd(cmd *cobra.Command, args []string, gflags *flags.GenerateFlags) error {
	serviceName := args[0]

	log.Debug().Str("service", serviceName).Str("zone", gflags.Zone).Msg("starting token generation")

	if gflags.Zone == "" {
		return config.ErrZoneRequired
	}

	req := cloudflare.TokenGenerationRequest{
		ServiceName: serviceName,
		ZoneName:    gflags.Zone,
		ExpiresOn:   gflags.ExpiresOn,
		Name:        gflags.Name,
	}

	if gflags.DryRun {
		return dryRun(cmd, gflags, req)
	}

	apiKey, err := resolveAPIKey(gflags.Token)
	if err != nil {
		return err
	}

	client, err := cloudflare.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("initialize Cloudflare client: %w", err)
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(gflags.Timeout)*time.Second)
	defer cancel()

	log.Debug().Dur("timeout", time.Duration(gflags.Timeout)*time.Second).Msg("generating token")

	result, err := cloudflare.GenerateToken(ctx, client, req)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	return outputToken(cmd, gflags, result)
}

// dryRun validates inputs and logs what would be created without making an API call.
//
// Parameters:
//   - cmd: Cobra command used for output.
//   - gflags: Generate flags controlling JSON output and output file.
//   - req: The token generation request containing service, zone, and name info.
//
// Returns:
//   - error: Non-nil if JSON output fails or inputs are invalid.
func dryRun(cmd *cobra.Command, gflags *flags.GenerateFlags, req cloudflare.TokenGenerationRequest) error {
	if !cloudflare.ServiceNamePattern.MatchString(req.ServiceName) {
		return fmt.Errorf("invalid service name: %s", req.ServiceName)
	}

	if !config.ZonePattern.MatchString(req.ZoneName) {
		return fmt.Errorf("invalid zone name: %s", req.ZoneName)
	}

	if req.ExpiresOn != "" {
		_, err := time.Parse(time.RFC3339, req.ExpiresOn)
		if err != nil {
			return fmt.Errorf("invalid expires-on %q: %w", req.ExpiresOn, err)
		}
	}

	tokenName := req.ServiceName + "." + req.ZoneName
	if req.Name != "" {
		tokenName = req.Name
	}

	if gflags.JSON {
		result := &cloudflare.TokenGenerationResult{
			Token:     "<dry-run-token>",
			Name:      tokenName,
			Zone:      req.ZoneName,
			ExpiresOn: req.ExpiresOn,
		}

		return outputTokenJSON(cmd, result)
	}

	outputPath := gflags.Output
	if outputPath != "" {
		logging.Printf("Dry run: would write token to file: %s", outputPath)
	} else {
		logging.Printf("Dry run: would write token to stdout")
	}

	logging.Printf("Dry run: would create token %s for zone %s", tokenName, req.ZoneName)

	return nil
}
