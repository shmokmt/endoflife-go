package endoflife

import (
	"context"
	"fmt"
	"net/http"
)

// GetTags retrieves a list of tags.
func (c *Client) GetTags(ctx context.Context) (*URIListResponse, error) {
	var result URIListResponse
	if err := c.doRequest(ctx, http.MethodGet, "/tags", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTagProducts retrieves a list of products with a specific tag.
func (c *Client) GetTagProducts(ctx context.Context, tagName string) (*ProductListResponse, error) {
	if tagName == "" {
		return nil, fmt.Errorf("tag name is required")
	}

	path := fmt.Sprintf("/tags/%s", tagName)
	var result ProductListResponse
	if err := c.doRequest(ctx, http.MethodGet, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
