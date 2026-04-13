package servercard

// Options configures the Server Card generation.
// Fields that can't be read from the go-sdk Server (which has unexported fields)
// must be provided here.
type Options struct {
	// Name in reverse-DNS format (e.g., "io.github.olgasafonova/gleif-mcp-server").
	// Required.
	Name string

	// Version string, should follow semver (e.g., "1.4.0").
	// Required.
	Version string

	// Description of the server's functionality.
	// Required.
	Description string

	// Title is an optional human-readable display name.
	Title string

	// WebsiteURL links to the server's homepage or documentation.
	WebsiteURL string

	// Repository describes the source code location.
	Repository *Repository

	// Icons for display in client UIs.
	Icons []Icon

	// Remotes describes HTTP-based transport endpoints.
	Remotes []Remote

	// Meta holds additional metadata.
	Meta map[string]any

	// Provider describes who operates this server.
	Provider *Provider
}

// Provider describes the server operator (not part of SEP-2127 schema,
// but useful metadata stored in _meta).
type Provider struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}
