// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/rs/zerolog/log"
)

const (
	// productionBaseURL is the Cloudflare v4 API endpoint. It is applied explicitly because
	// cloudflare.NewClient prepends DefaultClientOptions(), which maps the ambient
	// CLOUDFLARE_BASE_URL environment variable onto the client's base URL. Without an
	// explicit pin, a single environment variable could redirect authenticated requests,
	// and the bearer token they carry, to an arbitrary host.
	productionBaseURL = "https://api.cloudflare.com/client/v4/"

	// responseHeaderTimeout bounds the wait for response headers. The Cloudflare SDK
	// defaults to ten minutes, which is far too long for an interactive CLI whose commands
	// already apply their own deadline through the request context.
	responseHeaderTimeout = time.Minute

	// maxRedirects bounds how many redirects a single request may follow.
	maxRedirects = 10
)

// ambientSDKEnvVars are the environment variables cloudflare.NewClient reads on its own.
// This tool authenticates with an API token and pins the API endpoint, so every one of
// them is ignored. NewClient warns when any is set so that a user who exported one is not
// left believing it took effect.
var ambientSDKEnvVars = []string{
	"CLOUDFLARE_API_TOKEN",
	"CLOUDFLARE_API_KEY",
	"CLOUDFLARE_API_USER_SERVICE_KEY",
	"CLOUDFLARE_EMAIL",
	"CLOUDFLARE_BASE_URL",
	"CLOUDFLARE_CUSTOM_HEADERS",
}

// ambientAuthHeaders are the request headers cloudflare.NewClient derives from the ambient
// environment. This tool authenticates with a bearer token only, so they are removed rather
// than forwarding credentials the caller never supplied to this command.
var ambientAuthHeaders = []string{
	"X-Auth-Key",
	"X-Auth-Email",
	"X-Auth-User-Service-Key",
}

// newHTTPClient builds the HTTP client used for every Cloudflare API call.
//
// The transport is a clone of http.DefaultTransport so proxy handling, TLS verification and
// connection pooling keep their standard behaviour, with a bounded wait for response
// headers. Redirects that change host are refused outright, because an API response has no
// legitimate reason to move a request to a different origin.
//
// No client-level Timeout is set. Each command derives its own deadline from the request
// context, and that deadline is what the user controls through the --timeout flag; a
// client-level ceiling here would silently cap it.
//
// Returns:
//   - *http.Client: A client configured for the Cloudflare API.
func newHTTPClient() *http.Client {
	return &http.Client{
		Transport:     newTransport(),
		CheckRedirect: refuseCrossHostRedirect,
	}
}

// newTransport clones the default transport and applies the response header bound.
//
// Returns:
//   - http.RoundTripper: The transport to use, or the package default when it cannot be
//     cloned.
func newTransport() http.RoundTripper {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return http.DefaultTransport
	}

	bounded := transport.Clone()
	bounded.ResponseHeaderTimeout = responseHeaderTimeout

	return bounded
}

// refuseCrossHostRedirect rejects a redirect that would move the request to a different
// host, and caps the length of the redirect chain.
//
// Parameters:
//   - req: The redirect target.
//   - via: The requests already attempted, oldest first.
//
// Returns:
//   - error: ErrTooManyRedirects once the chain is too long, ErrCrossHostRedirect when the
//     host changes, otherwise nil to allow the redirect.
func refuseCrossHostRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return fmt.Errorf("%w: stopped after %d redirects", ErrTooManyRedirects, maxRedirects)
	}

	if len(via) > 0 && !strings.EqualFold(req.URL.Host, via[0].URL.Host) {
		return fmt.Errorf("%w: %s to %s", ErrCrossHostRedirect, via[0].URL.Host, req.URL.Host)
	}

	return nil
}

// clientOptions builds the option list applied to the Cloudflare SDK client.
//
// The order matters. cloudflare.NewClient places DefaultClientOptions() ahead of this
// list, so every option here is applied last and therefore wins: the HTTP client replaces
// the SDK's, the API token replaces any CLOUDFLARE_API_TOKEN, and the pinned base URL
// replaces any CLOUDFLARE_BASE_URL. The header deletions come last of all so they remove
// the legacy auth headers those environment-derived options installed.
//
// Parameters:
//   - apiToken: The Cloudflare API token used for authentication.
//   - httpClient: The HTTP client used to send requests.
//
// Returns:
//   - []option.RequestOption: The options to construct the client with.
func clientOptions(apiToken string, httpClient *http.Client) []option.RequestOption {
	baseOpts := []option.RequestOption{
		option.WithHTTPClient(httpClient),
		option.WithAPIToken(apiToken),
		// WithBaseURL assigns the explicit base URL, which outranks the SDK's default base
		// URL, so it must be applied after every option that could set either one.
		option.WithBaseURL(productionBaseURL),
	}

	opts := make([]option.RequestOption, 0, len(baseOpts)+len(ambientAuthHeaders))
	opts = append(opts, baseOpts...)

	for _, header := range ambientAuthHeaders {
		opts = append(opts, option.WithHeaderDel(header))
	}

	return opts
}

// warnIgnoredEnvironment logs a warning naming any Cloudflare SDK environment variable that
// is set but has no effect on this tool.
func warnIgnoredEnvironment() {
	var ignored []string

	for _, name := range ambientSDKEnvVars {
		if _, set := os.LookupEnv(name); set {
			ignored = append(ignored, name)
		}
	}

	if len(ignored) == 0 {
		return
	}

	log.Warn().
		Strs("variables", ignored).
		Str("apiEndpoint", productionBaseURL).
		Msg("ignoring Cloudflare SDK environment variables: the API endpoint is pinned and authentication uses the resolved API token")
}
