# endoflife-go

Go client library for the [endoflife.date](https://endoflife.date/) API v1.

## Installation

```bash
go get github.com/shmokmt/endoflife-go
```

## CLI

### Installation

```bash
go install github.com/shmokmt/endoflife-go/cmd/endoflife@latest
```

### Commands

- `endoflife products` - List all products
- `endoflife product <name>` - Get product details
  - `--release <version>` - Get specific release info
  - `--latest` - Get latest release info
- `endoflife version` - Show version

### Options

- `--json` - Output in JSON format
- `--timeout <duration>` - HTTP timeout (default: 30s)

### Examples

Check the latest release for Apache Airflow in JSON format:

```bash
endoflife product apache-airflow --latest --json
```

```json
{
  "schema_version": "1.2.0",
  "generated_at": "2026-01-06T00:32:36+00:00",
  "result": {
    "name": "3",
    "label": "3",
    "releaseDate": "2024-09-24",
    "isLts": false,
    "isEoas": false,
    "isEol": false,
    "isDiscontinued": false,
    "isMaintained": true,
    "latest": {
      "name": "3.1.5",
      "date": "2026-02-07",
      "link": "https://airflow.apache.org/docs/apache-airflow/stable/release_notes.html#airflow-3-0-1-2025-01-27"
    }
  }
}
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/shmokmt/endoflife-go"
)

func main() {
    client := endoflife.NewClient()
    ctx := context.Background()

    // Get list of all products
    products, err := client.GetProducts(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, product := range products.Result {
        fmt.Printf("%s (%s)\n", product.Label, product.Name)
    }
}
```

### Get Product Details

```go
product, err := client.GetProduct(ctx, "python")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Product: %s\n", product.Result.Label)
fmt.Printf("Category: %s\n", product.Result.Category)
fmt.Printf("Releases: %d\n", len(product.Result.Releases))
```

### Check EOL Status

```go
release, err := client.GetRelease(ctx, "python", "3.12")
if err != nil {
    log.Fatal(err)
}

if release.Result.IsEOL {
    fmt.Printf("Python 3.12 is EOL since %s\n", release.Result.EOLFrom.String())
} else {
    fmt.Println("Python 3.12 is still maintained")
}
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

client := endoflife.NewClientWithOptions(
    endoflife.WithHTTPClient(httpClient),
    endoflife.WithUserAgent("my-app/1.0"),
)
```

### Error Handling

```go
product, err := client.GetProduct(ctx, "nonexistent")
if err != nil {
    if endoflife.IsNotFound(err) {
        fmt.Println("Product not found")
    } else if endoflife.IsRateLimited(err) {
        fmt.Println("Rate limit exceeded")
    } else {
        log.Fatal(err)
    }
}
```

## API Methods

| Method | Description |
|--------|-------------|
| `GetIndex(ctx)` | Get API index |
| `GetProducts(ctx)` | Get all products (summary) |
| `GetProductsFull(ctx)` | Get all products (full details) |
| `GetProduct(ctx, name)` | Get a specific product |
| `GetRelease(ctx, product, release)` | Get a specific release |
| `GetLatestRelease(ctx, product)` | Get the latest release |
| `GetCategories(ctx)` | Get all categories |
| `GetCategoryProducts(ctx, category)` | Get products in a category |
| `GetTags(ctx)` | Get all tags |
| `GetTagProducts(ctx, tag)` | Get products with a tag |
| `GetIdentifiers(ctx)` | Get all identifier types |
| `GetIdentifierDetails(ctx, type)` | Get identifier details |

## License

MIT
