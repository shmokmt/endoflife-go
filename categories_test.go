package endoflife

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetCategories(t *testing.T) {
	expected := URIListResponse{
		SchemaVersion: "1.2.0",
		Total:         3,
		Result: []URI{
			{Name: "lang", URI: "/api/v1/categories/lang"},
			{Name: "os", URI: "/api/v1/categories/os"},
			{Name: "framework", URI: "/api/v1/categories/framework"},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetCategories(ctx)
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

func TestGetCategoryProducts(t *testing.T) {
	expected := ProductListResponse{
		SchemaVersion: "1.2.0",
		Total:         2,
		Result: []ProductSummary{
			{Name: "python", Label: "Python", Category: "lang"},
			{Name: "go", Label: "Go", Category: "lang"},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/categories/lang" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetCategoryProducts(ctx, "lang")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, result.Total)
	}
}

func TestGetCategoryProducts_EmptyName(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetCategoryProducts(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty category name")
	}
}
