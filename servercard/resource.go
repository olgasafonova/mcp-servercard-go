package servercard

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ResourceURI is the MCP resource URI for the Server Card.
// SEP-2127 states the card should also be available as an MCP resource.
// Note: the exact URI is not yet normative in the spec text; this value
// follows the convention from the PR summary.
const ResourceURI = "mcp://server-card.json"

// RegisterResource adds the Server Card as an MCP resource on the server.
// Clients that are already connected can read the card via resources/read
// without making an HTTP request.
func RegisterResource(server *mcp.Server, card *ServerCard) {
	data, err := card.JSON()
	if err != nil {
		panic("servercard: failed to marshal card: " + err.Error())
	}

	server.AddResource(&mcp.Resource{
		URI:         ResourceURI,
		Name:        "MCP Server Card",
		Description: "SEP-2127 Server Card describing this server's identity, transports, and capabilities.",
		MIMEType:    "application/json",
	}, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      ResourceURI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})
}
