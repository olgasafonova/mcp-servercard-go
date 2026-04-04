package servercard

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Attach builds a Server Card from opts, registers it as an MCP resource on
// the server, and returns an http.Handler to serve it at WellKnownPath.
//
// Typical usage with a standard http.ServeMux:
//
//	mux := http.NewServeMux()
//	mux.Handle("/mcp", mcpHandler)
//	mux.Handle(servercard.WellKnownPath, servercard.Attach(server, opts))
func Attach(server *mcp.Server, opts Options) http.Handler {
	card := Build(opts)
	RegisterResource(server, card)
	return Handler(card)
}
