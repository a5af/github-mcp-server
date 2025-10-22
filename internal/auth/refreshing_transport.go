package auth

import (
	"bytes"
	
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	// MaxRequestBodySize is the maximum request body size we'll clone for retry (10MB)
	// Requests larger than this won't be retried if they fail with 401
	MaxRequestBodySize = 10 * 1024 * 1024 // 10MB
)

// RefreshingAuthTransport wraps an HTTP transport with automatic token refresh on 401 errors.
//
// When a request receives a 401 Unauthorized response, this transport will:
// 1. Force an immediate token refresh from the auth provider
// 2. Retry the request once with the fresh token
// 3. Return the retry response (even if it's another 401)
//
// Retry Protection:
// - Only retries once per request (prevents infinite loops)
// - Request bodies larger than MaxRequestBodySize (10MB) cannot be retried
// - If retry also receives 401, returns it without further retry
//
// Thread Safety:
// This transport is safe for concurrent use across multiple goroutines.
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
	return t.roundTripWithRetry(req, false)
}

// roundTripWithRetry executes the request with optional retry logic
func (t *RefreshingAuthTransport) roundTripWithRetry(req *http.Request, isRetry bool) (*http.Response, error) {
	// Clone the request so we can retry it if needed
	reqCopy, bodyTooLarge, err := cloneRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to clone request: %w", err)
	}

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
	if resp.StatusCode == http.StatusUnauthorized && !isRetry {
		requestInfo := fmt.Sprintf("%s %s", req.Method, req.URL.Path)

		if bodyTooLarge {
			fmt.Fprintf(os.Stderr, "[refreshing-auth-transport] Got 401 for %s, but request body too large (>10MB) to retry\n", requestInfo)
			return resp, nil
		}

		fmt.Fprintf(os.Stderr, "[refreshing-auth-transport] Got 401 for %s, forcing token refresh\n", requestInfo)

		// Drain and close the response body before retrying
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// Force token refresh by invalidating current token and calling GetToken
		t.mu.Lock()
		if err := t.provider.refreshToken(req.Context()); err != nil {
			t.mu.Unlock()
			return nil, fmt.Errorf("failed to refresh GitHub App token for %s: %w (check credentials and installation ID)", requestInfo, err)
		}
		_ , err := t.provider.GetToken()
		if err != nil {
			t.mu.Unlock()
			return nil, fmt.Errorf("failed to get new token after refresh for %s: %w", requestInfo, err)
		}
		t.mu.Unlock()

		fmt.Fprintf(os.Stderr, "[refreshing-auth-transport] Token refreshed successfully, retrying %s\n", requestInfo)

		// Retry with new token (pass isRetry=true to prevent infinite loop)
		retryResp, retryErr := t.roundTripWithRetry(req, true)

		// Log if retry also failed with 401
		if retryErr == nil && retryResp != nil && retryResp.StatusCode == http.StatusUnauthorized {
			fmt.Fprintf(os.Stderr, "[refreshing-auth-transport] WARNING: Retry for %s also received 401 - token may be invalid or permissions insufficient\n", requestInfo)
		}

		return retryResp, retryErr
	}

	return resp, nil
}

// cloneRequest creates a shallow copy of the request with a cloned body if present.
// Returns (clonedRequest, bodyTooLarge, error).
// If bodyTooLarge is true, the request body was larger than MaxRequestBodySize and
// the clone will have an empty body (making retry unsafe).
func cloneRequest(req *http.Request) (*http.Request, bool, error) {
	r := req.Clone(req.Context())

	// If there's a body, we need to preserve it for potential retry
	if req.Body != nil && req.Body != http.NoBody {
		// Read the body into memory with size limit
		bodyBytes, err := io.ReadAll(io.LimitReader(req.Body, MaxRequestBodySize+1))
		if err != nil {
			return nil, false, fmt.Errorf("failed to read request body: %w", err)
		}

		// Check if body was too large
		if len(bodyBytes) > MaxRequestBodySize {
			// Body too large - can't safely retry
			// Restore original body for the first attempt, but don't clone for retry
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes[:MaxRequestBodySize]))
			r.Body = req.Body
			return r, true, nil
		}

		// Restore the original body and set the clone's body
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return r, false, nil
}
