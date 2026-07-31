# mcp-servercard-go

![lint](https://github.com/olgasafonova/mcp-servercard-go/actions/workflows/lint.yml/badge.svg)
[![CodeScene Average Code Health](https://codescene.io/projects/83042/status-badges/average-code-health)](https://codescene.io/projects/83042)

Reference implementation of SEP-2127 MCP Server Cards for Go.

Go library implementing [SEP-2127 MCP Server Cards](https://github.com/modelcontextprotocol/modelcontextprotocol/pull/2127). Pure-stdlib HTTP handler; no MCP SDK dependency.

A Server Card is a JSON document served at `/.well-known/mcp-server-card` that describes an MCP server before connection: identity, transports, protocol versions, and connection guidance. This enables pre-connect discovery without a full initialization handshake.

> **Note on `$schema` URL:** The spec references `https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json`. As of v0.3.0, this URL returns 404 — schema hosting is pending the SEP merge. The library emits the URL the spec mandates; treat the schema as informational until hosted.

## Status under MCP `2026-07-28`: complement to `server/discover`, not a duplicate of it

Protocol revision `2026-07-28` removed the `initialize` handshake and made the in-band `server/discover` RPC a MUST (SEP-2575). That looks like it supersedes pre-connect Server Cards. It does not, and the reason is in the payloads.

Checked against `go-sdk@v1.7.0` source (`mcp/protocol.go`, `mcp/server.go:903`), a `DiscoverResult` carries exactly three things: `supportedVersions`, `capabilities`, and `instructions`. **It carries no server identity at all**: no name, version, description, icons, repository, or endpoint URL.

| Question a client has | Answered by |
|---|---|
| Where is this server's endpoint? | Server Card: `remotes[].url` |
| What auth headers must I send? | Server Card: `remotes[].headers` |
| Who publishes it, what is it, icon, repo, website? | Server Card |
| Which transport type does that endpoint use? | Server Card: `remotes[].type` |
| What capabilities does it expose? | `server/discover` |
| How am I meant to use it? | `server/discover`, via `instructions` |
| Which protocol versions? | **both**: `remotes[].supportedProtocolVersions` and `supportedVersions` |

Protocol versions are the one genuine overlap. Everything else is disjoint, and the ordering is what matters: `server/discover` is a JSON-RPC call to an endpoint, so a client must already know the URL and any required headers before it can make one. The Server Card is what supplies those. It sits **upstream** of in-band discovery rather than competing with it.

The practical split: use the Server Card for catalogue-shaped questions answered before a connection exists (registries, directories, install-time UX), and `server/discover` for what the live endpoint currently speaks. A client that has already connected should prefer `server/discover` for versions and capabilities, since it reflects the running server rather than a published document.

**Upstream tracking.** SEP-2127 is still an open PR, delegated to [`modelcontextprotocol/experimental-ext-server-card`](https://github.com/modelcontextprotocol/experimental-ext-server-card) under the Extensions Track. That repo is active (last push 27-07-2026; recent work on catalog-based discovery and ETag revalidation), and a core maintainer has favoured adopting `server.json` as-is as the Server Card. Treat this library as tracking a moving extension, not a frozen spec: the shape may change before it lands.

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
- **Breaking:** Removed `Capabilities`, `Requires`, and `Authentication` fields per SEP-2127 scope reduction (capabilities and auth belong in `initialize`, not in static discovery). *Since `2026-07-28` there is no `initialize`: capabilities moved to `server/discover`. The scope reduction still holds, and the split is now the one described under Status above.*
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
