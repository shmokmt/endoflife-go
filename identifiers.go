package endoflife

import (
	"context"
	"fmt"
	"net/http"
)

// GetIdentifiers retrieves a list of identifier types.
func (c *Client) GetIdentifiers(ctx context.Context) (*URIListResponse, error) {
	var result URIListResponse
	if err := c.doRequest(ctx, http.MethodGet, "/identifiers", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIdentifierDetails retrieves details for a specific identifier type.
func (c *Client) GetIdentifierDetails(ctx context.Context, identifierType string) (*IdentifierListResponse, error) {
	if identifierType == "" {
		return nil, fmt.Errorf("identifier type is required")
	}

	path := fmt.Sprintf("/identifiers/%s", identifierType)
	var result IdentifierListResponse
	if err := c.doRequest(ctx, http.MethodGet, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
