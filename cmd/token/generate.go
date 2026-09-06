// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/gogeneratecftoken/internal/cloudflare"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/config"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/credentials"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/flags"
	"github.com/nicholas-fedor/gogeneratecftoken/internal/logging"
)

const (
	minDurationLen = 2
	hoursPerDay    = 24
	//nolint:gosec // G101: placeholder token value for dry-run output, not a credential
	dryRunTokenPlaceholder = "<dry-run-token>"
)

var (
	errTTLExclusive         = errors.New("--ttl and --expires-on are mutually exclusive")
	errInvalidServiceName   = errors.New("invalid service name")
	errUnrecognizedDuration = errors.New("unrecognized duration")
	errDurationOverflow     = errors.New("duration overflows maximum")
	errNonPositiveDuration  = errors.New("duration must be greater than zero")
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

  # Generate with a TTL and account ID
  goGenerateCFToken token generate myapp --ttl 7d --account-id acc123

  # Output as JSON to stdout
  goGenerateCFToken token generate myapp --json

  # Write token to file with no stdout
  goGenerateCFToken token generate myapp --output /tmp/token.txt --format none`,
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

// NewDeprecatedGenerateCommand returns a root-level generate alias for v1
// compatibility. It is marked deprecated in favor of `token generate`.
//
// Returns:
//   - *cobra.Command: The deprecated generate command.
func NewDeprecatedGenerateCommand() *cobra.Command {
	cmd := newGenerateCommand()
	cmd.Deprecated = `use "token generate" instead`
	cmd.GroupID = ""

	return cmd
}

// loadConfig loads the configuration and sets the zone flag if not already provided.
//
// Parameters:
//   - cmd: Cobra command for retrieving the config flag value.
//   - gflags: Generate flags providing the optional zone. Zone may be updated.
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
		if !errors.Is(err, config.ErrZoneRequired) {
			return fmt.Errorf("load config: %w", err)
		}

		d := config.Default()
		cfg = &d
	}

	if gflags.Zone == "" {
		gflags.Zone = cmp.Or(cfg.Zone, os.Getenv(config.EnvVarZone))
		log.Debug().Str("zone", gflags.Zone).Msg("using zone from config or CF_ZONE")
	}

	if gflags.AccountID == "" {
		gflags.AccountID = cfg.AccountID
	}

	if gflags.Name == "" {
		gflags.Name = cfg.TokenName
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
//   - args: Positional arguments. args[0] is the service name.
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
		AccountID:   gflags.AccountID,
		ExpiresOn:   gflags.ExpiresOn,
		Name:        gflags.Name,
	}

	if gflags.TTL != "" {
		if req.ExpiresOn != "" {
			return errTTLExclusive
		}

		dur, err := parseDuration(gflags.TTL)
		if err != nil {
			return fmt.Errorf("invalid ttl %q: %w", gflags.TTL, err)
		}

		req.ExpiresOn = time.Now().UTC().Add(dur).Format(time.RFC3339)
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
		return fmt.Errorf("%w: %s", errInvalidServiceName, req.ServiceName)
	}

	if !config.ZonePattern.MatchString(req.ZoneName) {
		return fmt.Errorf("%w: %s", config.ErrInvalidZone, req.ZoneName)
	}

	if req.ExpiresOn != "" {
		_, err := time.Parse(time.RFC3339, req.ExpiresOn)
		if err != nil {
			return fmt.Errorf("invalid expires-on %q: %w", req.ExpiresOn, err)
		}
	}

	if gflags.TTL != "" {
		_, err := parseDuration(gflags.TTL)
		if err != nil {
			return fmt.Errorf("invalid ttl %q: %w", gflags.TTL, err)
		}
	}

	tokenName := req.ServiceName + "." + req.ZoneName
	if req.Name != "" {
		tokenName = req.Name
	}

	if gflags.JSON {
		result := &cloudflare.TokenGenerationResult{
			Token:     dryRunTokenPlaceholder,
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

// parseDuration converts a human-friendly duration string to a time.Duration.
// It supports Go's standard suffixes (ns, us/µs, ms, s, m, h) and adds
// "d" for days (1d = 24h).
//
// Parameters:
//   - value: The duration string (e.g. "24h", "7d").
//
// Returns:
//   - time.Duration: The parsed duration.
//   - error: Non-nil if the format is not recognized.
func parseDuration(value string) (time.Duration, error) {
	if len(value) < minDurationLen {
		return 0, fmt.Errorf("%w: %s", errUnrecognizedDuration, value)
	}

	var (
		dur time.Duration
		err error
	)

	if before, ok := strings.CutSuffix(value, "d"); ok {
		days, perr := strconv.Atoi(before)
		if perr != nil {
			return 0, fmt.Errorf("%w: %s", errUnrecognizedDuration, value)
		}

		const maxDays = int(math.MaxInt64 / int64(hoursPerDay*time.Hour))
		if days > maxDays {
			return 0, fmt.Errorf("%w: %s", errDurationOverflow, value)
		}

		dur = time.Duration(days) * hoursPerDay * time.Hour
	} else {
		dur, err = time.ParseDuration(value)
		if err != nil {
			return 0, fmt.Errorf("%w: %s: %w", errUnrecognizedDuration, value, err)
		}
	}

	if dur <= 0 {
		return 0, fmt.Errorf("%w: %s", errNonPositiveDuration, value)
	}

	return dur, nil
}
