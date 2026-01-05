package endoflife

import (
	"context"
	"fmt"
	"net/http"
)

// GetProducts retrieves a list of product summaries.
func (c *Client) GetProducts(ctx context.Context) (*ProductListResponse, error) {
	var result ProductListResponse
	if err := c.doRequest(ctx, http.MethodGet, "/products", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProductsFull retrieves a list of full product details.
func (c *Client) GetProductsFull(ctx context.Context) (*FullProductListResponse, error) {
	var result FullProductListResponse
	if err := c.doRequest(ctx, http.MethodGet, "/products/full", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProduct retrieves detailed information for a specific product.
func (c *Client) GetProduct(ctx context.Context, productName string) (*ProductResponse, error) {
	if productName == "" {
		return nil, fmt.Errorf("product name is required")
	}

	path := fmt.Sprintf("/products/%s", productName)
	var result ProductResponse
	if err := c.doRequest(ctx, http.MethodGet, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRelease retrieves release information for a specific product release.
func (c *Client) GetRelease(ctx context.Context, productName, releaseName string) (*ProductReleaseResponse, error) {
	if productName == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if releaseName == "" {
		return nil, fmt.Errorf("release name is required")
	}

	path := fmt.Sprintf("/products/%s/releases/%s", productName, releaseName)
	var result ProductReleaseResponse
	if err := c.doRequest(ctx, http.MethodGet, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLatestRelease retrieves the latest release information for a specific product.
func (c *Client) GetLatestRelease(ctx context.Context, productName string) (*ProductReleaseResponse, error) {
	if productName == "" {
		return nil, fmt.Errorf("product name is required")
	}

	path := fmt.Sprintf("/products/%s/releases/latest", productName)
	var result ProductReleaseResponse
	if err := c.doRequest(ctx, http.MethodGet, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
