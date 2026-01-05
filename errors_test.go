package endoflife

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiErr   *APIError
		expected string
	}{
		{
			name:     "without retry after",
			apiErr:   &APIError{StatusCode: 404, Message: "resource not found"},
			expected: "API error: 404 resource not found",
		},
		{
			name:     "with retry after",
			apiErr:   &APIError{StatusCode: 429, Message: "rate limit exceeded", RetryAfter: 60},
			expected: "API error: 429 rate limit exceeded (retry after 60 seconds)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiErr.Error()
			if result != tt.expected {
				t.Errorf("APIError.Error() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			expected: true,
		},
		{
			name:     "APIError with 404",
			err:      &APIError{StatusCode: 404, Message: "not found"},
			expected: true,
		},
		{
			name:     "APIError with 500",
			err:      &APIError{StatusCode: 500, Message: "server error"},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotFound() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrRateLimited",
			err:      ErrRateLimited,
			expected: true,
		},
		{
			name:     "APIError with 429",
			err:      &APIError{StatusCode: 429, Message: "rate limited", RetryAfter: 60},
			expected: true,
		},
		{
			name:     "APIError with 404",
			err:      &APIError{StatusCode: 404, Message: "not found"},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRateLimited(tt.err)
			if result != tt.expected {
				t.Errorf("IsRateLimited() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNotModified(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrNotModified",
			err:      ErrNotModified,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotModified(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotModified() = %v, want %v", result, tt.expected)
			}
		})
	}
}
