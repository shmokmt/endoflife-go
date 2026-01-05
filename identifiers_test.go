package endoflife

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetIdentifiers(t *testing.T) {
	expected := URIListResponse{
		SchemaVersion: "1.2.0",
		Total:         3,
		Result: []URI{
			{Name: "purl", URI: "/api/v1/identifiers/purl"},
			{Name: "repology", URI: "/api/v1/identifiers/repology"},
			{Name: "cpe", URI: "/api/v1/identifiers/cpe"},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/identifiers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetIdentifiers(ctx)
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

func TestGetIdentifierDetails(t *testing.T) {
	expected := IdentifierListResponse{
		SchemaVersion: "1.2.0",
		Total:         2,
		Result: []IdentifierMapping{
			{Identifier: "pkg:pypi/python", Product: URI{Name: "python", URI: "/api/v1/products/python"}},
			{Identifier: "pkg:golang/go", Product: URI{Name: "go", URI: "/api/v1/products/go"}},
		},
	}

	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/identifiers/purl" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetIdentifierDetails(ctx, "purl")
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

func TestGetIdentifierDetails_EmptyType(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	_, err := client.GetIdentifierDetails(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty identifier type")
	}
}
