# mcp-servercard-go

![lint](https://github.com/olgasafonova/mcp-servercard-go/actions/workflows/lint.yml/badge.svg)
[![CodeScene Average Code Health](https://codescene.io/projects/83042/status-badges/average-code-health)](https://codescene.io/projects/83042)

Reference implementation of SEP-2127 MCP Server Cards for Go.

Go library implementing [SEP-2127 MCP Server Cards](https://github.com/modelcontextprotocol/modelcontextprotocol/pull/2127). Pure-stdlib HTTP handler; no MCP SDK dependency.

A Server Card is a JSON document served at `/.well-known/mcp-server-card` that describes an MCP server before connection: identity, transports, protocol versions, and connection guidance. This enables pre-connect discovery without a full initialization handshake.

> **Note on `$schema` URL:** The spec references `https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json`. As of v0.3.0, this URL returns 404 — schema hosting is pending the SEP merge. The library emits the URL the spec mandates; treat the schema as informational until hosted.

## Usage

```go
import "github.com/olgasafonova/mcp-servercard-go/servercard"

cardHandler, err := servercard.Attach(servercard.Options{
    Name:        "io.github.olgasafonova/gleif-mcp-server",
    Version:     "1.4.0",
    Description: "Access the GLEIF LEI database for company verification.",
    Title:       "GLEIF MCP Server",
    WebsiteURL:  "https://github.com/olgasafonova/gleif-mcp-server",
    Remotes: []servercard.Remote{{
        Type: "streamable-http",
        URL:  "/mcp",
    }},
    Provider: &servercard.Provider{
        Name: "Olga Safonova",
        URL:  "https://github.com/olgasafonova",
    },
})
if err != nil {
    log.Fatal(err)
}

// Mount alongside your MCP handler.
mux := http.NewServeMux()
mux.Handle("/mcp", mcpHandler)
mux.Handle(servercard.WellKnownPath, cardHandler)
```

## What it does

1. Builds a Server Card JSON document conforming to the SEP-2127 schema
2. Serves it at `/.well-known/mcp-server-card` with correct CORS and caching headers

## What's new in v0.3.0

- **Breaking:** Removed `RegisterResource()`, `ResourceURI`, and the `mcp://server-card.json` MCP resource endpoint to align with SEP-2127 (PR #2443, March 2026): discovery is pre-connection only via `.well-known`.
- **Breaking:** `Attach()` no longer takes a `*mcp.Server` parameter — signature is now `Attach(opts) (http.Handler, error)`.
- **Breaking:** Removed `Capabilities`, `Requires`, and `Authentication` fields per SEP-2127 scope reduction (capabilities and auth belong in `initialize`, not in static discovery).
- Library no longer depends on `github.com/modelcontextprotocol/go-sdk`. Pure stdlib.

## What's new in v0.2.0

- `Auth.Schemes` normalization: empty or nil slices serialize consistently
- `WellKnownPathFor(card)` for multi-server setups: returns `/.well-known/mcp-server-card/{name}` so multiple servers on the same host each get a distinct discovery path

## API

| Function | Purpose |
|----------|---------|
| `Attach(opts)` | One-line setup: builds card and returns HTTP handler |
| `Build(opts)` | Builds a `*ServerCard` struct from options |
| `Handler(card)` | Returns an `http.Handler` serving the card JSON |
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
