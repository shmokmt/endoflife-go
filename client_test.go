package endoflife

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := NewClientWithOptions(WithBaseURL(server.URL))
	return client, server
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client.BaseURL != DefaultBaseURL {
		t.Errorf("expected BaseURL %s, got %s", DefaultBaseURL, client.BaseURL)
	}
	if client.HTTPClient == nil {
		t.Error("expected HTTPClient to be set")
	}
	if client.UserAgent != "endoflife-go/"+Version {
		t.Errorf("expected UserAgent endoflife-go/%s, got %s", Version, client.UserAgent)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	customBaseURL := "https://custom.example.com"
	customUserAgent := "custom-agent/1.0"
	customHTTPClient := &http.Client{Timeout: 60 * time.Second}

	client := NewClientWithOptions(
		WithBaseURL(customBaseURL),
		WithUserAgent(customUserAgent),
		WithHTTPClient(customHTTPClient),
	)

	if client.BaseURL != customBaseURL {
		t.Errorf("expected BaseURL %s, got %s", customBaseURL, client.BaseURL)
	}
	if client.UserAgent != customUserAgent {
		t.Errorf("expected UserAgent %s, got %s", customUserAgent, client.UserAgent)
	}
	if client.HTTPClient != customHTTPClient {
		t.Error("expected HTTPClient to be custom client")
	}
}

func TestGetIndex(t *testing.T) {
	expected := URIListResponse{
		SchemaVersion: "1.2.0",
		Total:         3,
		Result: []URI{
			{Name: "products", URI: "/api/v1/products"},
			{Name: "categories", URI: "/api/v1/categories"},
			{Name: "tags", URI: "/api/v1/tags"},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("unexpected Accept header: %s", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetIndex(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.SchemaVersion != expected.SchemaVersion {
		t.Errorf("expected schema_version %s, got %s", expected.SchemaVersion, result.SchemaVersion)
	}
	if result.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, result.Total)
	}
	if len(result.Result) != len(expected.Result) {
		t.Errorf("expected %d results, got %d", len(expected.Result), len(result.Result))
	}
}

func TestHandleHTTPError_NotFound(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.GetIndex(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsNotFound(err) {
		t.Errorf("expected NotFound error, got %v", err)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("expected APIError, got %T", err)
	} else if apiErr.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestHandleHTTPError_RateLimited(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.GetIndex(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsRateLimited(err) {
		t.Errorf("expected RateLimited error, got %v", err)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("expected APIError, got %T", err)
	} else {
		if apiErr.StatusCode != 429 {
			t.Errorf("expected status code 429, got %d", apiErr.StatusCode)
		}
		if apiErr.RetryAfter != 60 {
			t.Errorf("expected RetryAfter 60, got %d", apiErr.RetryAfter)
		}
	}
}

func TestHandleHTTPError_NotModified(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.GetIndex(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsNotModified(err) {
		t.Errorf("expected NotModified error, got %v", err)
	}
}

func TestContextCancellation(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.GetIndex(ctx)
	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}
