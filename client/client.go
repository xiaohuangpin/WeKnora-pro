// Package client provides the implementation for interacting with the WeKnora API
// This package encapsulates CRUD operations for server resources and provides a friendly interface for callers
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the client for interacting with the WeKnora service
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// ClientOption defines client configuration options
type ClientOption func(*Client)

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithToken sets the authentication token
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// NewClient creates a new client instance
func NewClient(baseURL string, options ...ClientOption) *Client {
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// doRequest executes an HTTP request
func (c *Client) doRequest(ctx context.Context,
	method, path string, body interface{}, query url.Values,
) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	if len(query) > 0 {
		url = fmt.Sprintf("%s?%s", url, query.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-API-Key", c.token)
	}
	if requestID := ctx.Value("RequestID"); requestID != nil {
		req.Header.Set("X-Request-ID", requestID.(string))
	}

	return c.httpClient.Do(req)
}

// parseResponse parses an HTTP response
func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	if target == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
