package endoflife

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the endoflife.date API.
	DefaultBaseURL = "https://endoflife.date/api/v1"

	// DefaultTimeout is the default timeout for HTTP requests.
	DefaultTimeout = 30 * time.Second

	// Version is the library version.
	Version = "0.1.0"
)

// Client is the client for the endoflife.date API.
type Client struct {
	// BaseURL is the base URL for the API.
	BaseURL string

	// HTTPClient is a customizable HTTP client.
	HTTPClient *http.Client

	// UserAgent is the User-Agent to set on requests.
	UserAgent string
}

// NewClient creates a new Client with default settings.
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		UserAgent: "endoflife-go/" + Version,
	}
}

// Option is a function type for configuring the client.
type Option func(*Client)

// WithBaseURL sets the base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithUserAgent sets the User-Agent.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// NewClientWithOptions creates a new Client with the given options.
func NewClientWithOptions(opts ...Option) *Client {
	c := NewClient()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doRequest executes an HTTP request and processes the response.
func (c *Client) doRequest(ctx context.Context, method, path string, result any) error {
	reqURL, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.handleHTTPError(resp); err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// handleHTTPError handles HTTP error responses.
func (c *Client) handleHTTPError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusMovedPermanently:
		return nil
	case http.StatusNotModified:
		return ErrNotModified
	case http.StatusNotFound:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "resource not found",
		}
	case http.StatusTooManyRequests:
		retryAfter := 0
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "rate limit exceeded",
			RetryAfter: retryAfter,
		}
	default:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}
}

// GetIndex retrieves the API index (list of main endpoints).
func (c *Client) GetIndex(ctx context.Context) (*URIListResponse, error) {
	var result URIListResponse
	if err := c.doRequest(ctx, http.MethodGet, "/", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
