package endoflife_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shmokmt/endoflife-go"
)

func Example() {
	client := endoflife.NewClient()
	ctx := context.Background()

	// Get list of products
	products, err := client.GetProducts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total products: %d\n", products.Total)
}

func Example_getProduct() {
	client := endoflife.NewClient()
	ctx := context.Background()

	// Get Python product details
	product, err := client.GetProduct(ctx, "python")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Product: %s\n", product.Result.Label)
	fmt.Printf("Category: %s\n", product.Result.Category)
	fmt.Printf("Releases: %d\n", len(product.Result.Releases))
}

func Example_checkEOL() {
	client := endoflife.NewClient()
	ctx := context.Background()

	// Get Python 3.12 release information
	release, err := client.GetRelease(ctx, "python", "3.12")
	if err != nil {
		log.Fatal(err)
	}

	if release.Result.IsEOL {
		fmt.Printf("Python 3.12 is EOL since %s\n", release.Result.EOLFrom.String())
	} else {
		fmt.Println("Python 3.12 is still maintained")
	}
}

func Example_customHTTPClient() {
	// Use a custom HTTP client
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
		},
	}

	client := endoflife.NewClientWithOptions(
		endoflife.WithHTTPClient(httpClient),
		endoflife.WithUserAgent("my-app/1.0"),
	)

	ctx := context.Background()
	_, _ = client.GetProducts(ctx)
}

func Example_errorHandling() {
	client := endoflife.NewClient()
	ctx := context.Background()

	_, err := client.GetProduct(ctx, "nonexistent-product")
	if err != nil {
		if endoflife.IsNotFound(err) {
			fmt.Println("Product not found")
		} else if endoflife.IsRateLimited(err) {
			fmt.Println("Rate limit exceeded, please retry later")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func Example_getCategories() {
	client := endoflife.NewClient()
	ctx := context.Background()

	// Get all categories
	categories, err := client.GetCategories(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, cat := range categories.Result {
		fmt.Printf("Category: %s\n", cat.Name)
	}
}

func Example_getTags() {
	client := endoflife.NewClient()
	ctx := context.Background()

	// Get all tags
	tags, err := client.GetTags(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, tag := range tags.Result {
		fmt.Printf("Tag: %s\n", tag.Name)
	}
}
