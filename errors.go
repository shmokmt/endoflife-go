package endoflife

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a resource is not found (404).
	ErrNotFound = errors.New("resource not found")

	// ErrRateLimited is returned when the rate limit is exceeded (429).
	ErrRateLimited = errors.New("rate limit exceeded")

	// ErrNotModified is returned when the resource has not been modified (304).
	ErrNotModified = errors.New("resource not modified")
)

// APIError represents an error response from the API.
type APIError struct {
	StatusCode int
	Message    string
	RetryAfter int // Retry wait time in seconds for 429 responses
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("API error: %d %s (retry after %d seconds)",
			e.StatusCode, e.Message, e.RetryAfter)
	}
	return fmt.Sprintf("API error: %d %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether the error is a 404 error.
func IsNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return false
}

// IsRateLimited reports whether the error is a 429 error.
func IsRateLimited(err error) bool {
	if errors.Is(err, ErrRateLimited) {
		return true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 429
	}
	return false
}

// IsNotModified reports whether the error is a 304 error.
func IsNotModified(err error) bool {
	return errors.Is(err, ErrNotModified)
}
