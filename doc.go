// Package endoflife provides a Go client for the endoflife.date API v1.
//
// The endoflife.date API provides end-of-life information for various
// software products, including operating systems, programming languages,
// frameworks, and more.
//
// Basic usage:
//
//	client := endoflife.NewClient()
//	ctx := context.Background()
//
//	products, err := client.GetProducts(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, product := range products.Result {
//	    fmt.Println(product.Name)
//	}
//
// Custom configuration:
//
//	client := endoflife.NewClientWithOptions(
//	    endoflife.WithHTTPClient(customHTTPClient),
//	    endoflife.WithUserAgent("my-app/1.0"),
//	)
//
// Error handling:
//
//	product, err := client.GetProduct(ctx, "python")
//	if endoflife.IsNotFound(err) {
//	    // Handle 404 error
//	}
//	if endoflife.IsRateLimited(err) {
//	    // Handle rate limiting
//	}
//
// For more information about the API, see https://endoflife.date/docs/api
package endoflife
