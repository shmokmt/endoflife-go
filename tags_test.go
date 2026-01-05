package endoflife

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetTags(t *testing.T) {
	expected := URIListResponse{
		SchemaVersion: "1.2.0",
		Total:         2,
		Result: []URI{
			{Name: "programming", URI: "/api/v1/tags/programming"},
			{Name: "database", URI: "/api/v1/tags/database"},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tags" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetTags(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, result.Total)
	}
	if len(result.Result) != len(expected.Result) {
		t.Errorf("expected %d results, got %d", len(expected.Result), len(result.Result))
	}
}

func TestGetTagProducts(t *testing.T) {
	expected := ProductListResponse{
		SchemaVersion: "1.2.0",
		Total:         2,
		Result: []ProductSummary{
			{Name: "python", Label: "Python", Tags: []string{"programming"}},
			{Name: "go", Label: "Go", Tags: []string{"programming"}},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tags/programming" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetTagProducts(ctx, "programming")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, result.Total)
	}
}

func TestGetTagProducts_EmptyName(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetTagProducts(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty tag name")
	}
}
