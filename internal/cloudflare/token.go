// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

//nolint:wrapcheck
package cloudflare

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/shared"
	"github.com/cloudflare/cloudflare-go/v7/user"
	"github.com/cloudflare/cloudflare-go/v7/zones"
)

const (
	// zoneReadPermission is the Cloudflare permission group ID for zone read access.
	zoneReadPermission = "c8fed203ed3043cba015a93ad1616f1f"
	// dnsWritePermission is the Cloudflare permission group ID for DNS write access.
	dnsWritePermission = "4755a26eedb94da69e1066d98aa820be"
)

// ServiceNamePattern validates that service names contain only alphanumeric characters, hyphens, and underscores.
var ServiceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// TokenInfo holds display information for an API token.
type TokenInfo struct {
	// ID is the unique token identifier.
	ID string `json:"id"`
	// Name is the human-readable token name.
	Name string `json:"name"`
	// Status indicates the current token status (e.g., active, disabled).
	Status string `json:"status"`
}

// TokenGenerationRequest holds parameters for generating a token.
type TokenGenerationRequest struct {
	// ServiceName is the service identifier used in the default token name.
	ServiceName string
	// ZoneName is the DNS zone for which the token is granted permissions.
	ZoneName string
	// AccountID is the optional Cloudflare account ID for zone disambiguation.
	AccountID string
	// ExpiresOn is the optional expiration time in RFC3339 format.
	ExpiresOn string
	// Name overrides the default token name (service.zone).
	Name string
}

// TokenGenerationResult holds the result of token generation.
type TokenGenerationResult struct {
	// Token is the generated API token value.
	Token string
	// Name is the token name that was assigned.
	Name string
	// Zone is the DNS zone the token is scoped to.
	Zone string
	// ExpiresOn is the expiration time in RFC3339 format, if set.
	ExpiresOn string
}

// resolveZoneID looks up the zone ID for the given zone name.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - api: The API implementation to use for zone lookup.
//   - zoneName: The DNS zone name to resolve.
//   - accountID: Optional account ID to disambiguate zones across accounts.
//
// Returns:
//   - string: The resolved zone ID.
//   - error: Non-nil if the zone is not found, multiple zones match, or the API call fails.
func resolveZoneID(ctx context.Context, api API, zoneName, accountID string) (string, error) {
	params := zones.ZoneListParams{Name: cloudflare.F(zoneName)}

	if accountID != "" {
		params.Account = cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(accountID)})
	}

	zones, err := api.ListZones(ctx, params)
	if err != nil {
		return "", err
	}

	switch len(zones) {
	case 0:
		return "", fmt.Errorf("%w: %s", ErrZoneNotFound, zoneName)
	case 1:
		return zones[0].ID, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrMultipleZonesFound, zoneName)
	}
}

// buildTokenPolicy constructs the IAM policy for zone read and DNS write permissions.
//
// Parameters:
//   - zoneID: The zone ID to scope the policy to.
//
// Returns:
//   - []shared.TokenPolicyParam: The policy parameters for token creation.
func buildTokenPolicy(zoneID string) []shared.TokenPolicyParam {
	permissions := []shared.TokenPolicyPermissionGroupParam{
		{ID: cloudflare.F(zoneReadPermission)},
		{ID: cloudflare.F(dnsWritePermission)},
	}

	resources := shared.TokenPolicyResourcesIAMResourcesTypeObjectStringParam{
		"com.cloudflare.api.account.zone." + zoneID: "*",
	}

	return []shared.TokenPolicyParam{{
		Effect:           cloudflare.F(shared.TokenPolicyEffectAllow),
		PermissionGroups: cloudflare.F(permissions),
		Resources:        cloudflare.F(shared.TokenPolicyResourcesUnionParam(resources)),
	}}
}

// GenerateToken creates a Cloudflare API token with DNS edit permissions.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - api: The API implementation to use for token creation.
//   - req: The token generation request containing service name, zone, and optional settings.
//
// Returns:
//   - *TokenGenerationResult: The generated token and its metadata.
//   - error: Non-nil if the service name is invalid, zone cannot be resolved, or creation fails.
func GenerateToken(ctx context.Context, api API, req TokenGenerationRequest) (*TokenGenerationResult, error) {
	if !ServiceNamePattern.MatchString(req.ServiceName) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidServiceName, req.ServiceName)
	}

	zID, err := resolveZoneID(ctx, api, req.ZoneName, req.AccountID)
	if err != nil {
		return nil, err
	}

	tokenName := req.ServiceName + "." + req.ZoneName
	if req.Name != "" {
		tokenName = req.Name
	}

	params := user.TokenNewParams{
		Name:     cloudflare.F(tokenName),
		Policies: cloudflare.F(buildTokenPolicy(zID)),
	}

	if req.ExpiresOn != "" {
		expiresAt, err := time.Parse(time.RFC3339, req.ExpiresOn)
		if err != nil {
			return nil, fmt.Errorf("invalid expires-on %q: %w", req.ExpiresOn, err)
		}

		params.ExpiresOn = cloudflare.F(expiresAt)
	}

	token, err := api.CreateAPIToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return &TokenGenerationResult{
		Token:     token.Value,
		Name:      tokenName,
		Zone:      req.ZoneName,
		ExpiresOn: req.ExpiresOn,
	}, nil
}
