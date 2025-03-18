package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

// APIInterface defines the methods required by GenerateToken and related functions.
type APIInterface interface {
	ListZones(ctx context.Context, zone ...string) ([]cloudflare.Zone, error)
	CreateAPIToken(ctx context.Context, token cloudflare.APIToken) (cloudflare.APIToken, error)
}
