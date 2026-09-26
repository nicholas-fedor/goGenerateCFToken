// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package cloudflare

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyTokenListPayload is a well-formed response whose empty result terminates the SDK's
// auto-pager immediately, so a single request is recorded per client call.
const emptyTokenListPayload = `{"result":[],"success":true,"errors":[],"messages":[]}`

// recordingTransport captures the requests the SDK attempts and answers each with a canned
// payload, so client construction can be asserted without touching the network.
type recordingTransport struct {
	requests []*http.Request
}

// RoundTrip records the request and replies with the canned payload.
//
// Parameters:
//   - req: The request the SDK attempted.
//
// Returns:
//   - *http.Response: A synthetic 200 response.
//   - error: Always nil.
func (rec *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec.requests = append(rec.requests, req.Clone(req.Context()))

	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(emptyTokenListPayload)),
		Request:    req,
	}, nil
}

// newRecordingClient builds a client whose requests are captured rather than sent.
//
// Parameters:
//   - apiToken: The Cloudflare API token used for authentication.
//
// Returns:
//   - *Client: The client under test.
//   - *recordingTransport: The transport recording its requests.
func newRecordingClient(apiToken string) (*Client, *recordingTransport) {
	recorder := &recordingTransport{}

	return newClient(apiToken, &http.Client{Transport: recorder}), recorder
}

func TestNewClient_PinsProductionEndpoint(t *testing.T) {
	t.Setenv("CLOUDFLARE_BASE_URL", "https://attacker.example.com/")

	client, recorder := newRecordingClient("test-token-12345")

	_, err := client.ListTokens(t.Context())
	require.NoError(t, err)
	require.Len(t, recorder.requests, 1)

	got := recorder.requests[0].URL
	assert.Equal(t, "api.cloudflare.com", got.Host)
	assert.Equal(t, "/client/v4/user/tokens", got.Path)
}

func TestNewClient_IgnoresAmbientAuthHeaders(t *testing.T) {
	t.Setenv("CLOUDFLARE_API_KEY", "ambient-global-key")
	t.Setenv("CLOUDFLARE_EMAIL", "ambient@example.com")
	t.Setenv("CLOUDFLARE_API_USER_SERVICE_KEY", "ambient-service-key")

	client, recorder := newRecordingClient("test-token-12345")

	_, err := client.ListTokens(t.Context())
	require.NoError(t, err)
	require.Len(t, recorder.requests, 1)

	header := recorder.requests[0].Header
	assert.Equal(t, "Bearer test-token-12345", header.Get("Authorization"))
	assert.Empty(t, header.Get("X-Auth-Key"))
	assert.Empty(t, header.Get("X-Auth-Email"))
	assert.Empty(t, header.Get("X-Auth-User-Service-Key"))
}

func TestNewClient_IgnoresAmbientAPIToken(t *testing.T) {
	t.Setenv("CLOUDFLARE_API_TOKEN", "ambient-token")

	client, recorder := newRecordingClient("explicit-token")

	_, err := client.ListTokens(t.Context())
	require.NoError(t, err)
	require.Len(t, recorder.requests, 1)

	assert.Equal(t, "Bearer explicit-token", recorder.requests[0].Header.Get("Authorization"))
}

func TestNewClient_UsesSuppliedHTTPClient(t *testing.T) {
	recorder := &recordingTransport{}
	httpClient := &http.Client{Transport: recorder}

	client := newClient("test-token-12345", httpClient)

	_, err := client.ListTokens(t.Context())
	require.NoError(t, err)
	assert.Len(t, recorder.requests, 1)
}

func TestNewHTTPClient(t *testing.T) {
	client := newHTTPClient()

	assert.NotNil(t, client.Transport)
	assert.NotNil(t, client.CheckRedirect)
}

func TestRefuseCrossHostRedirect(t *testing.T) {
	first, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"https://api.cloudflare.com/client/v4/user/tokens",
		nil,
	)
	require.NoError(t, err)

	sameHost, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"https://api.cloudflare.com/client/v4/user/tokens/abc",
		nil,
	)
	require.NoError(t, err)

	assert.NoError(t, refuseCrossHostRedirect(sameHost, []*http.Request{first}))

	otherHost, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"https://attacker.example.com/collect",
		nil,
	)
	require.NoError(t, err)

	err = refuseCrossHostRedirect(otherHost, []*http.Request{first})
	require.ErrorIs(t, err, ErrCrossHostRedirect)
}

func TestRefuseCrossHostRedirect_TooManyRedirects(t *testing.T) {
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"https://api.cloudflare.com/client/v4/user/tokens",
		nil,
	)
	require.NoError(t, err)

	chain := make([]*http.Request, maxRedirects)
	for i := range chain {
		chain[i] = req
	}

	err = refuseCrossHostRedirect(req, chain)
	require.ErrorIs(t, err, ErrTooManyRedirects)
}

func TestAmbientSDKEnvVarsCoverDocumentedSDKVariables(t *testing.T) {
	documented := []string{
		"CLOUDFLARE_API_TOKEN",
		"CLOUDFLARE_API_KEY",
		"CLOUDFLARE_API_USER_SERVICE_KEY",
		"CLOUDFLARE_EMAIL",
		"CLOUDFLARE_BASE_URL",
		"CLOUDFLARE_CUSTOM_HEADERS",
	}

	for _, name := range documented {
		assert.Contains(t, ambientSDKEnvVars, name)
	}
}
