// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
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

	// envCustomHeaders names the SDK environment variable that injects arbitrary request
	// headers. Its contents are parsed and removed rather than trusted.
	envCustomHeaders = "CLOUDFLARE_CUSTOM_HEADERS"

	// contentTypeHeader and acceptHeader are set by the SDK during request setup, before
	// any client option is applied.
	contentTypeHeader = "content-type"
	acceptHeader      = "accept"

	// jsonMediaType is the value the SDK uses for both of those headers.
	jsonMediaType = "application/json"
)

// protectedHeaders are headers the SDK owns and that must never be treated as ambient.
// Authorization carries the token this client exists to present, Content-Type is required
// on requests with a body, and Accept selects the response representation.
var protectedHeaders = []string{
	"authorization",
	contentTypeHeader,
	acceptHeader,
}

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
	envCustomHeaders,
}

// ambientAuthHeaders are the request headers cloudflare.NewClient derives from the ambient
// environment. This tool authenticates with a bearer token only, so they are removed rather
// than forwarding credentials the caller never supplied to this command.
var ambientAuthHeaders = []string{
	"X-Auth-Key",
	"X-Auth-Email",
	"X-Auth-User-Service-Key",
}

// ambientCustomHeaders returns the header names cloudflare.NewClient derives from
// CLOUDFLARE_CUSTOM_HEADERS, parsed with the same rules the SDK uses. Protected names are
// included so that the caller can tell an injected header from an absent one; whether a
// name is removed or restored is the caller's decision.
//
// Returns:
//   - []string: The injected header names, empty when the variable is unset.
func ambientCustomHeaders() []string {
	raw, set := os.LookupEnv(envCustomHeaders)
	if !set {
		return nil
	}

	names := make([]string, 0, strings.Count(raw, "\n")+1)

	for line := range strings.SplitSeq(raw, "\n") {
		before, _, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		if name := strings.TrimSpace(before); name != "" {
			names = append(names, name)
		}
	}

	return names
}

// isProtectedHeader reports whether a header name is one the SDK owns and must keep.
//
// Parameters:
//   - name: The candidate header name.
//
// Returns:
//   - bool: True when the name matches a protected header, ignoring case.
func isProtectedHeader(name string) bool {
	return slices.ContainsFunc(protectedHeaders, func(protected string) bool {
		return strings.EqualFold(name, protected)
	})
}

// newHTTPClient builds the HTTP client used for every Cloudflare API call.
//
// The transport is a clone of http.DefaultTransport so proxy handling, TLS verification and
// connection pooling keep their standard behaviour, with a bounded wait for response
// headers. Redirects that change origin are refused outright, because an API response has no
// legitimate reason to move a request to a different origin, and a scheme downgrade would
// otherwise put the bearer token on the wire in clear text.
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
		CheckRedirect: refuseCrossOriginRedirect,
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

// sameOrigin reports whether two URLs address the same origin. Scheme and host are both
// compared case-insensitively, as both are case-insensitive per RFC 3986.
//
// Parameters:
//   - target: The URL being redirected to.
//   - previous: The URL of the request being redirected from.
//
// Returns:
//   - bool: True when the two share a scheme and host.
func sameOrigin(target, previous *url.URL) bool {
	return strings.EqualFold(target.Scheme, previous.Scheme) &&
		strings.EqualFold(target.Host, previous.Host)
}

// refuseCrossOriginRedirect rejects a redirect that would move the request to a different
// origin, and caps the length of the redirect chain. A host change would move the bearer
// token to a third party, and a scheme change on the same host would expose it in clear
// text, so both are refused.
//
// Parameters:
//   - req: The redirect target.
//   - via: The requests already attempted, oldest first.
//
// Returns:
//   - error: ErrTooManyRedirects once the chain is too long, ErrCrossOriginRedirect when
//     the origin changes, otherwise nil to allow the redirect.
func refuseCrossOriginRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return fmt.Errorf("%w: stopped after %d redirects", ErrTooManyRedirects, maxRedirects)
	}

	if len(via) > 0 && !sameOrigin(req.URL, via[0].URL) {
		return fmt.Errorf("%w: %s to %s", ErrCrossOriginRedirect, via[0].URL, req.URL)
	}

	return nil
}

// clientOptions builds the option list applied to the Cloudflare SDK client.
//
// The order matters. cloudflare.NewClient places DefaultClientOptions() ahead of this
// list, so every option here is applied last and therefore wins: the HTTP client replaces
// the SDK's, the API token replaces any CLOUDFLARE_API_TOKEN, and the pinned base URL
// replaces any CLOUDFLARE_BASE_URL. The header deletions come after those so they remove
// what the environment injected, and the SDK-owned values are restored last because the
// SDK applies every option after its own request setup.
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

	removed := ambientCustomHeaders()
	injected := make(map[string]struct{}, len(removed))

	opts := make([]option.RequestOption, 0, len(baseOpts)+len(ambientAuthHeaders)+len(removed))
	opts = append(opts, baseOpts...)

	for _, header := range ambientAuthHeaders {
		opts = append(opts, option.WithHeaderDel(header))
	}

	for _, header := range removed {
		injected[strings.ToLower(header)] = struct{}{}

		// A protected name is left in place, because deleting it would strip something
		// the request needs. It is restored below instead.
		if !isProtectedHeader(header) {
			opts = append(opts, option.WithHeaderDel(header))
		}
	}

	// The SDK sets these two during request setup, before options are applied, so an
	// injected header of the same name has already replaced them and deleting it is not
	// enough. Put the SDK's value back, but only when one was actually targeted, so
	// requests that never carried the header are left as they were.
	if _, hit := injected[contentTypeHeader]; hit {
		opts = append(opts, option.WithHeader(contentTypeHeader, jsonMediaType))
	}

	if _, hit := injected[acceptHeader]; hit {
		opts = append(opts, option.WithHeader(acceptHeader, jsonMediaType))
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
