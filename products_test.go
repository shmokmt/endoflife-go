package endoflife

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetProducts(t *testing.T) {
	expected := ProductListResponse{
		SchemaVersion: "1.2.0",
		Total:         2,
		Result: []ProductSummary{
			{Name: "python", Label: "Python", Category: "lang", Tags: []string{"programming"}, URI: "/api/v1/products/python"},
			{Name: "go", Label: "Go", Category: "lang", Tags: []string{"programming"}, URI: "/api/v1/products/go"},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetProducts(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, result.Total)
	}
	if len(result.Result) != len(expected.Result) {
		t.Errorf("expected %d results, got %d", len(expected.Result), len(result.Result))
	}
	if result.Result[0].Name != "python" {
		t.Errorf("expected first product name 'python', got %s", result.Result[0].Name)
	}
}

func TestGetProductsFull(t *testing.T) {
	expected := FullProductListResponse{
		SchemaVersion: "1.2.0",
		Total:         1,
		Result: []ProductDetails{
			{
				Name:     "python",
				Label:    "Python",
				Category: "lang",
				Tags:     []string{"programming"},
				URI:      "/api/v1/products/python",
				Releases: []ProductRelease{},
			},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/full" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetProductsFull(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, result.Total)
	}
}

func TestGetProduct(t *testing.T) {
	expected := ProductResponse{
		SchemaVersion: "1.2.0",
		LastModified:  "2024-01-01T00:00:00Z",
		Result: ProductDetails{
			Name:     "python",
			Label:    "Python",
			Category: "lang",
			Tags:     []string{"programming"},
			URI:      "/api/v1/products/python",
			Releases: []ProductRelease{},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/python" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetProduct(ctx, "python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Result.Name != "python" {
		t.Errorf("expected product name 'python', got %s", result.Result.Name)
	}
}

func TestGetProduct_EmptyName(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetProduct(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty product name")
	}
}

func TestGetRelease(t *testing.T) {
	expected := ProductReleaseResponse{
		SchemaVersion: "1.2.0",
		LastModified:  "2024-01-01T00:00:00Z",
		Result: ProductRelease{
			Name:         "3.12",
			Label:        "3.12",
			IsLTS:        false,
			IsEOL:        false,
			IsMaintained: true,
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/python/releases/3.12" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetRelease(ctx, "python", "3.12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Result.Name != "3.12" {
		t.Errorf("expected release name '3.12', got %s", result.Result.Name)
	}
}

func TestGetRelease_EmptyProductName(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetRelease(ctx, "", "3.12")
	if err == nil {
		t.Fatal("expected error for empty product name")
	}
}

func TestGetRelease_EmptyReleaseName(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetRelease(ctx, "python", "")
	if err == nil {
		t.Fatal("expected error for empty release name")
	}
}

func TestGetLatestRelease(t *testing.T) {
	expected := ProductReleaseResponse{
		SchemaVersion: "1.2.0",
		Result: ProductRelease{
			Name:         "3.13",
			Label:        "3.13",
			IsMaintained: true,
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/python/releases/latest" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetLatestRelease(ctx, "python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Result.Name != "3.13" {
		t.Errorf("expected release name '3.13', got %s", result.Result.Name)
	}
}

func TestGetLatestRelease_EmptyName(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetLatestRelease(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty product name")
	}
}
