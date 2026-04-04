// Package servercard implements SEP-2127 MCP Server Cards.
//
// A Server Card is a JSON document that describes an MCP server before
// connection, enabling pre-connect discovery of capabilities, transports,
// and authentication requirements.
//
// See https://github.com/modelcontextprotocol/modelcontextprotocol/pull/2127
package servercard

import "encoding/json"

// SchemaURL is the JSON Schema URI for SEP-2127 Server Cards.
const SchemaURL = "https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json"

// ServerCard is the top-level SEP-2127 Server Card document.
type ServerCard struct {
	Schema       string         `json:"$schema"`
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Description  string         `json:"description"`
	Title        string         `json:"title,omitempty"`
	WebsiteURL   string         `json:"websiteUrl,omitempty"`
	Repository   *Repository    `json:"repository,omitempty"`
	Icons        []Icon         `json:"icons,omitempty"`
	Remotes      []Remote       `json:"remotes,omitempty"`
	Capabilities *Capabilities  `json:"capabilities,omitempty"`
	Requires     *Requires      `json:"requires,omitempty"`
	Meta         map[string]any `json:"_meta,omitempty"`
}

// Repository describes the source code location for the server.
type Repository struct {
	URL       string `json:"url"`
	Source    string `json:"source,omitempty"`
	Subfolder string `json:"subfolder,omitempty"`
}

// Icon represents a server icon for display in client UIs.
type Icon struct {
	Source   string   `json:"src"`
	MIMEType string   `json:"mimeType,omitempty"`
	Sizes    []string `json:"sizes,omitempty"`
	Theme    string   `json:"theme,omitempty"`
}

// Remote describes an HTTP-based transport endpoint.
type Remote struct {
	Type                      string   `json:"type"`
	URL                       string   `json:"url"`
	SupportedProtocolVersions []string `json:"supportedProtocolVersions,omitempty"`
	Headers                   []Header `json:"headers,omitempty"`
	Authentication            *Auth    `json:"authentication,omitempty"`
}

// Header describes a custom HTTP header the server expects.
type Header struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	IsRequired  bool     `json:"isRequired,omitempty"`
	IsSecret    bool     `json:"isSecret,omitempty"`
	Default     string   `json:"default,omitempty"`
	Choices     []string `json:"choices,omitempty"`
}

// Auth describes authentication requirements for a remote.
type Auth struct {
	Required bool     `json:"required"`
	Schemes  []string `json:"schemes"`
}

// Capabilities mirrors the MCP ServerCapabilities shape.
type Capabilities struct {
	Experimental map[string]any  `json:"experimental,omitempty"`
	Logging      *LoggingCap     `json:"logging,omitempty"`
	Completions  *CompletionsCap `json:"completions,omitempty"`
	Prompts      *PromptsCap     `json:"prompts,omitempty"`
	Resources    *ResourcesCap   `json:"resources,omitempty"`
	Tools        *ToolsCap       `json:"tools,omitempty"`
}

// LoggingCap indicates log message support.
type LoggingCap struct{}

// CompletionsCap indicates argument autocompletion support.
type CompletionsCap struct{}

// PromptsCap describes prompt template capabilities.
type PromptsCap struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCap describes resource capabilities.
type ResourcesCap struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// ToolsCap describes tool capabilities.
type ToolsCap struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Requires describes what the server needs from the client.
type Requires struct {
	Experimental map[string]any `json:"experimental,omitempty"`
	Sampling     *struct{}      `json:"sampling,omitempty"`
	Roots        *struct{}      `json:"roots,omitempty"`
	Elicitation  *struct{}      `json:"elicitation,omitempty"`
}

// JSON returns the Server Card as indented JSON bytes.
func (c *ServerCard) JSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
