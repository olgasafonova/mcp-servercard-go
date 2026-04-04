package servercard

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Attach builds a Server Card from opts, registers it as an MCP resource on
// the server, and returns an http.Handler to serve it at WellKnownPath.
//
// Returns an error if opts fails validation (missing required fields or
// invalid name format).
//
// Typical usage with a standard http.ServeMux:
//
//	mux := http.NewServeMux()
//	mux.Handle("/mcp", mcpHandler)
//	handler, err := servercard.Attach(server, opts)
//	mux.Handle(servercard.WellKnownPath, handler)
func Attach(server *mcp.Server, opts Options) (http.Handler, error) {
	card, err := Build(opts)
	if err != nil {
		return nil, err
	}
	RegisterResource(server, card)
	return Handler(card), nil
}
