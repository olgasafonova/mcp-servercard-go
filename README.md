# mcp-servercard-go

Reference implementation of SEP-2127 MCP Server Cards for Go.

Go library implementing [SEP-2127 MCP Server Cards](https://github.com/modelcontextprotocol/modelcontextprotocol/pull/2127). Drop-in middleware for any [go-sdk](https://github.com/modelcontextprotocol/go-sdk) MCP server.

A Server Card is a JSON document served at `/.well-known/mcp-server-card` that describes an MCP server before connection: identity, transports, capabilities, and authentication requirements. This enables pre-connect discovery without a full initialization handshake.

## Usage

```go
import "github.com/olgasafonova/mcp-servercard-go/servercard"

server := mcp.NewServer(&mcp.Implementation{
    Name:    "gleif-mcp-server",
    Version: "1.4.0",
}, nil)

// ... register tools, prompts, resources ...

// Attach returns an http.Handler for the well-known endpoint
// and registers the card as an MCP resource (resource registration
// may be removed in a future release per PR #2443).
cardHandler := servercard.Attach(server, servercard.Options{
    Name:        "io.github.olgasafonova/gleif-mcp-server",
    Version:     "1.4.0",
    Description: "Access the GLEIF LEI database for company verification.",
    Title:       "GLEIF MCP Server",
    WebsiteURL:  "https://github.com/olgasafonova/gleif-mcp-server",
    Remotes: []servercard.Remote{{
        Type: "streamable-http",
        URL:  "/mcp",
        Authentication: &servercard.Auth{Required: false, Schemes: []string{}},
    }},
    Capabilities: &servercard.Capabilities{
        Tools:   &servercard.ToolsCap{ListChanged: false},
        Prompts: &servercard.PromptsCap{ListChanged: false},
    },
    Provider: &servercard.Provider{
        Name: "Olga Safonova",
        URL:  "https://github.com/olgasafonova",
    },
})

// Mount alongside your MCP handler.
mux := http.NewServeMux()
mux.Handle("/mcp", mcpHandler)
mux.Handle(servercard.WellKnownPath, cardHandler)
```

## What it does

1. Builds a Server Card JSON document conforming to the SEP-2127 schema
2. Serves it at `/.well-known/mcp-server-card` with correct CORS and caching headers
3. Registers it as an MCP resource at `mcp://server-card.json` so connected clients can read it too. Note: PR #2443 (merged March 26) removed this endpoint from the spec; this feature may be deprecated in a future release.

## What's new in v0.2.0

- `Auth.Schemes` normalization: empty or nil slices serialize consistently
- `WellKnownPathFor(card)` for multi-server setups: returns `/.well-known/mcp-server-card/{name}` so multiple servers on the same host each get a distinct discovery path

## API

| Function | Purpose |
|----------|---------|
| `Attach(server, opts)` | One-line setup: builds card, registers resource, returns HTTP handler |
| `Build(opts)` | Builds a `*ServerCard` struct from options |
| `Handler(card)` | Returns an `http.Handler` serving the card JSON |
| `RegisterResource(server, card)` | Registers the card as an MCP resource |
| `WellKnownPathFor(card)` | Returns `/.well-known/mcp-server-card/{name}` for multi-server sub-paths |

## Integrated in

- [miro-mcp-server](https://github.com/olgasafonova/miro-mcp-server) (91 tools, 5 prompts)
- [mediawiki-mcp-server](https://github.com/olgasafonova/mediawiki-mcp-server) (40+ tools)

## References

- [SEP-2127: MCP Server Cards](https://github.com/modelcontextprotocol/modelcontextprotocol/pull/2127)
- [MCP go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- [Skills Over MCP Interest Group](https://github.com/modelcontextprotocol/experimental-ext-skills)

## License

MIT
