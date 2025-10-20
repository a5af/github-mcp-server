package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// RefreshingAuthTransport wraps an HTTP transport with automatic token refresh on 401 errors
type RefreshingAuthTransport struct {
	transport http.RoundTripper
	provider  *GitHubAppAuthProvider
	mu        sync.Mutex
}

// NewRefreshingAuthTransport creates a new transport that automatically refreshes tokens on 401 errors
func NewRefreshingAuthTransport(transport http.RoundTripper, provider *GitHubAppAuthProvider) *RefreshingAuthTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &RefreshingAuthTransport{
		transport: transport,
		provider:  provider,
	}
}

// RoundTrip executes a single HTTP transaction, with automatic retry on 401 errors
func (t *RefreshingAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request so we can retry it if needed
	reqCopy := cloneRequest(req)

	// Get current token and set Authorization header
	token, err := t.provider.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	reqCopy.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	resp, err := t.transport.RoundTrip(reqCopy)
	if err != nil {
		return resp, err
	}

	// If we get a 401, the token may have expired - force refresh and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Fprintf(os.Stderr, "[refreshing-auth-transport] Got 401 response, forcing token refresh\n")

		// Drain and close the response body before retrying
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// Force token refresh by invalidating current token and calling GetToken
		t.mu.Lock()
		if err := t.provider.refreshToken(req.Context()); err != nil {
			t.mu.Unlock()
			return nil, fmt.Errorf("failed to refresh token after 401: %w", err)
		}
		newToken, err := t.provider.GetToken()
		if err != nil {
			t.mu.Unlock()
			return nil, fmt.Errorf("failed to get new token after refresh: %w", err)
		}
		t.mu.Unlock()

		fmt.Fprintf(os.Stderr, "[refreshing-auth-transport] Token refreshed, retrying request\n")

		// Clone the original request again for retry
		retryReq := cloneRequest(req)
		retryReq.Header.Set("Authorization", "Bearer "+newToken)

		// Retry with new token
		return t.transport.RoundTrip(retryReq)
	}

	return resp, nil
}

// cloneRequest creates a shallow copy of the request with a cloned body if present
func cloneRequest(req *http.Request) *http.Request {
	r := req.Clone(req.Context())

	// If there's a body, we need to preserve it for potential retry
	if req.Body != nil && req.Body != http.NoBody {
		// Read the body into memory
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			// If we can't read the body, just use an empty body
			r.Body = http.NoBody
		} else {
			// Restore the original body and set the clone's body
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	return r
}
