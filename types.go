package endoflife

import (
	"encoding/json"
	"time"
)

// Date is a custom type for date (YYYY-MM-DD format).
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Date.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil || *s == "" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON implements json.Marshaler for Date.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(d.Format("2006-01-02"))
}

// String returns the date in YYYY-MM-DD format.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02")
}

// URI represents an API endpoint reference.
type URI struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// Identifier represents a product identifier.
type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ProductVersion represents the latest version information of a product.
type ProductVersion struct {
	Name string  `json:"name"`
	Date *Date   `json:"date,omitempty"`
	Link *string `json:"link,omitempty"`
}

// ProductRelease represents an individual release of a product.
type ProductRelease struct {
	Name             string          `json:"name"`
	Codename         *string         `json:"codename,omitempty"`
	Label            string          `json:"label"`
	ReleaseDate      Date            `json:"releaseDate"`
	IsLTS            bool            `json:"isLts"`
	LTSFrom          *Date           `json:"ltsFrom,omitempty"`
	IsEOAS           bool            `json:"isEoas"`
	EOASFrom         *Date           `json:"eoasFrom,omitempty"`
	IsEOL            bool            `json:"isEol"`
	EOLFrom          *Date           `json:"eolFrom,omitempty"`
	IsDiscontinued   bool            `json:"isDiscontinued"`
	DiscontinuedFrom *Date           `json:"discontinuedFrom,omitempty"`
	IsEOES           *bool           `json:"isEoes,omitempty"`
	EOESFrom         *Date           `json:"eoesFrom,omitempty"`
	IsMaintained     bool            `json:"isMaintained"`
	Latest           *ProductVersion `json:"latest,omitempty"`
	Custom           map[string]any  `json:"custom,omitempty"`
}

// ProductSummary represents summary information of a product.
type ProductSummary struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Aliases  []string `json:"aliases"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	URI      string   `json:"uri"`
}

// ProductLabels represents label information of a product.
type ProductLabels struct {
	EOAS         *string `json:"eoas,omitempty"`
	Discontinued *string `json:"discontinued,omitempty"`
	EOL          string  `json:"eol"`
	EOES         *string `json:"eoes,omitempty"`
}

// ProductLinks represents link information of a product.
type ProductLinks struct {
	Icon          *string `json:"icon,omitempty"`
	HTML          string  `json:"html"`
	ReleasePolicy *string `json:"releasePolicy,omitempty"`
}

// ProductDetails represents detailed information of a product.
type ProductDetails struct {
	Name           string           `json:"name"`
	Label          string           `json:"label"`
	Aliases        []string         `json:"aliases"`
	Category       string           `json:"category"`
	Tags           []string         `json:"tags"`
	URI            string           `json:"uri"`
	VersionCommand *string          `json:"versionCommand,omitempty"`
	Identifiers    []Identifier     `json:"identifiers"`
	Labels         ProductLabels    `json:"labels"`
	Links          ProductLinks     `json:"links"`
	Releases       []ProductRelease `json:"releases"`
}

// URIListResponse represents a response containing a list of URIs.
type URIListResponse struct {
	SchemaVersion string `json:"schema_version"`
	GeneratedAt   string `json:"generated_at,omitempty"`
	Total         int    `json:"total"`
	Result        []URI  `json:"result"`
}

// ProductListResponse represents a response containing a list of product summaries.
type ProductListResponse struct {
	SchemaVersion string           `json:"schema_version"`
	GeneratedAt   string           `json:"generated_at,omitempty"`
	Total         int              `json:"total"`
	Result        []ProductSummary `json:"result"`
}

// FullProductListResponse represents a response containing a list of full product details.
type FullProductListResponse struct {
	SchemaVersion string           `json:"schema_version"`
	GeneratedAt   string           `json:"generated_at,omitempty"`
	Total         int              `json:"total"`
	Result        []ProductDetails `json:"result"`
}

// ProductResponse represents a response containing product details.
type ProductResponse struct {
	SchemaVersion string         `json:"schema_version"`
	GeneratedAt   string         `json:"generated_at,omitempty"`
	LastModified  string         `json:"last_modified"`
	Result        ProductDetails `json:"result"`
}

// ProductReleaseResponse represents a response containing release details.
type ProductReleaseResponse struct {
	SchemaVersion string         `json:"schema_version"`
	GeneratedAt   string         `json:"generated_at,omitempty"`
	LastModified  string         `json:"last_modified,omitempty"`
	Result        ProductRelease `json:"result"`
}

// IdentifierMapping represents a mapping between an identifier and a product.
type IdentifierMapping struct {
	Identifier string `json:"identifier"`
	Product    URI    `json:"product"`
}

// IdentifierListResponse represents a response containing a list of identifier mappings.
type IdentifierListResponse struct {
	SchemaVersion string              `json:"schema_version"`
	GeneratedAt   string              `json:"generated_at,omitempty"`
	Total         int                 `json:"total"`
	Result        []IdentifierMapping `json:"result"`
}
