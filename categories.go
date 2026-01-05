package endoflife

import (
	"context"
	"fmt"
	"net/http"
)

// GetCategories retrieves a list of categories.
func (c *Client) GetCategories(ctx context.Context) (*URIListResponse, error) {
	var result URIListResponse
	if err := c.doRequest(ctx, http.MethodGet, "/categories", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCategoryProducts retrieves a list of products in a specific category.
func (c *Client) GetCategoryProducts(ctx context.Context, categoryName string) (*ProductListResponse, error) {
	if categoryName == "" {
		return nil, fmt.Errorf("category name is required")
	}

	path := fmt.Sprintf("/categories/%s", categoryName)
	var result ProductListResponse
	if err := c.doRequest(ctx, http.MethodGet, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
