# SEP-2127 ecosystem check

Date of check: 23-07-2026. Scope: bead `mcp-servercard-go-4fe`. Library state at time of check: v0.3.0, last commit 04-05-2026.

Every factual claim below carries a confidence tag per `rules/confidence-scoring.md`.

## Headline

**The watch target did not stall. It moved, and it moved through the two things this library hard-codes: the track and the path.**

The bead's premise (SEP stalled in Draft, maintainer nags unanswered, watch the extension repo instead) was correct on 21-07-2026 for the *SEP PR as a venue*, but wrong about the proposal. Between 08-06-2026 and 20-07-2026 the working group did three consequential things:

1. Refactored SEP-2127 from **Standards Track to Extensions Track**, slimming the SEP to a charter and delegating the wire format to `experimental-ext-server-card`. That refactor PR is **merged**. `[HIGH: PR #2893 state MERGED at 2026-06-26T15:52:43Z via gh api, and PR #2127 carries the "extension" label]`
2. **Dropped `.well-known` as the Server Card location.** The recommended location is now `GET <streamable-http-url>/server-card`, with any unreserved URI allowed. `[HIGH: normative text in the current docs/discovery.md "Hosted Server Card Location", plus the PR #22 that landed it, closed 2026-06-08]`
3. **Replaced the MCP Catalog with the external AI Catalog** as the discovery entrypoint, at `/.well-known/ai-catalog.json`. `[HIGH: current docs/discovery.md "AI Catalog" section; PR #42 "Use AI Catalog for Server Card discovery" and issue #26 "Remove the MCP Catalog" both closed 2026-07-20]`

**Consequence for this repo: `servercard.WellKnownPath = "/.well-known/mcp-server-card"` and `WellKnownPathFor()` now serve a path the spec explicitly considered and rejected.** The rejection is written down: discovery.md's "Alternatives considered" lists a `.well-known` URI first, reasoning that `.well-known` is for site-wide metadata while a server card is application-level metadata. `[HIGH: verbatim in docs/discovery.md, fetched raw from the repo]`

SEP-2127 itself is still **open, not merged, `reviewDecision: CHANGES_REQUESTED`**, created 21-01-2026, last updated 20-07-2026. `[HIGH: gh pr view on the PR]` So the status is not "merged" and not "dead". It is "the SEP text is now a charter, and the real spec lives elsewhere and is moving weekly."

## 1. Production feedback

**b-gutman / gateway.pipeworx.io is real and still running, and has already migrated off the path this library serves.** `[HIGH: live curl 23-07-2026 plus their own PR comments]`

- `https://gateway.pipeworx.io/.well-known/mcp-server-card` returns **301 to `https://gateway.pipeworx.io/mcp/server-card`**. `[HIGH: live curl, HTTP 301 with that Location]`
- `https://gateway.pipeworx.io/.well-known/mcp/catalog.json` returns **200**. `[HIGH: live curl]`
- `https://gateway.pipeworx.io/.well-known/ai-catalog.json` returns **405**, so they have not yet moved to the AI Catalog well-known that the current spec text recommends. `[HIGH: live curl]`

Their 09-06-2026 comment states the migration explicitly: cards now at `<streamable-http-url>/server-card`, discovery via a catalog, v1 `.well-known/mcp-server-card[/...]` paths 301 to v2 equivalents so early adopters keep working. Pack count grew from 755 to 810 between 05-06 and 09-06-2026. `[HIGH: their PR #2127 comments, cross-confirmed by the live 301 above]`

The bead's two design validations from their 05-06 report still hold and are worth keeping:

- The server.json-subset shape required no new metadata sourcing. `[MEDIUM: single implementer's report, not corroborated by a second implementer]`
- Path collision with `.well-known/oauth-*` was a non-issue. `[MEDIUM: same single report; now moot anyway since the card left .well-known]`

Their "which of my 810 do I pick" default-pointer gap remains unresolved in-card. The catalog `identifier` plus the AI Catalog `metadata` surface is where the WG pushed it. `[MEDIUM: inferred from issue #41's title and the discovery.md entry table; no explicit WG statement found saying "this is the answer to the default-pointer problem"]`

Beyond pipeworx, **I could not find a second production deployment.** No search surfaced another live `/server-card` or `/.well-known/mcp-server-card` endpoint from an independent operator. `[LOW: absence of evidence from a bounded GitHub-side search, not a systematic crawl. Treat as "none found", not "none exists".]`

## 2. Interop: wong2/mcp-server-detector

**The detector still targets the old path and looks unmaintained.**

- Last push **07-05-2026**, 2 stars, not archived. `[HIGH: gh api repos/wong2/mcp-server-detector]`
- Its README states it checks the active site for a card at `/.well-known/mcp-server-card` and follows PR #2127. `[HIGH: README fetched raw]`

So the interop answer is a shrug with a short shelf life: this library and that extension agree on both the path and the card shape, which means they interoperate today by both being wrong about the current spec. Testing against it would validate a v1 contract that the spec has since replaced. `[MEDIUM: reasoning from two HIGH facts (its target path, the spec's new path); no live interop test was run]`

**Recommendation: do not spend the interop test. It validates the obsolete surface.** Retest only if the detector picks up `/server-card` and AI Catalog discovery.

## 3. Scope risk: YE-YI7's value-metadata question

**Answered by construction, not by the authors, and the answer is "out of card".** `[MEDIUM: inferred from WG artifacts below; I found no direct maintainer reply to YE-YI7 in the PR thread]`

- YE-YI7 opened `experimental-ext-server-card#41`, "Reference implementation showcase: value/selection metadata riding AI Catalog `metadata` + Server Card `_meta`". It was **closed 2026-07-20**. `[HIGH: gh api issue list]`
- The `_meta` extension surface was documented as a deliberate rider pattern in issue #27, closed 13-06-2026. `[HIGH: gh api issue list]`
- The current `schema.ts` has **no pricing, SLA, quality, or payment fields**, and its doc comment says a Server Card describes "only what is needed to discover and connect: identity, transport, and protocol versions". `[HIGH: schema.ts fetched raw, read in full]`

**Schema-addition risk from the value layer is currently LOW.** The pressure was routed to `_meta` and to the AI Catalog `metadata` field, both of which this library already supports (the `Meta map[string]any` field) or does not need to. `[MEDIUM: judgment call built on the HIGH facts above]`

The live scope risk is elsewhere. `experimental-ext-server-card#30`, "Add optional tool metadata to Server Card for offline discovery", is **open** since 17-06-2026. That is a direct attempt to reintroduce primitives, which the current schema and discovery.md explicitly exclude. If it lands it is a schema addition here. `[HIGH: issue open per gh api; exclusion is verbatim in schema.ts and discovery.md]`

## 4. The competitive fact the bead does not contain

**There is now an official Go SDK implementation of Server Cards in flight.**

`modelcontextprotocol/go-sdk#1024`, "mcp/experimental: add ServerCard convenience helper (SEP-2127)", by SamMorrowDrums. Draft, open, created 25-06-2026, last updated 21-07-2026, 1480 additions across 10 files including a full `mcp/experimental/servercard` package with tests and docs. `[HIGH: gh pr view and the PR file list]`

Its API surface overlaps this library almost exactly, and diverges where the spec moved: `[HIGH: read from the PR's patch for servercard.go]`

| Concept | go-sdk#1024 | mcp-servercard-go v0.3.0 |
|---|---|---|
| Media type | `MediaType = "application/mcp-server-card+json"` | `Content-Type: application/json` |
| Path | `DefaultPath = "/server-card"` (relative to the streamable-HTTP URL) | `WellKnownPath = "/.well-known/mcp-server-card"` |
| Build | `BuildServerCard(impl, opts...)` with functional options | `Build(opts)` / `Attach(opts)` with a struct |
| Serve | `Handler(card)`, `Mount(mux, path, card)` | `Handler(card)`, `Attach` returns a handler |
| Validation | `(*ServerCard).Validate()` | validation inside `Build` |
| ETag / 304 | strong ETag plus `If-None-Match` handling | none |

On 20-07-2026 tadasant said on issue #16 that the SDK reference PRs will be marked ready for review and merged pending SDK maintainer approval, but that **merging will not block SEP review**. `[HIGH: comment quoted from gh api on issue #16]`

This matters for positioning. The README calls this repo the "Reference implementation of SEP-2127 MCP Server Cards for Go". Once #1024 merges, the official reference implementation for Go is in the official SDK. `[HIGH: the PR exists and is scoped as the Go reference implementation, per its own body and issue #16]`

## 5. Everything else that drifted since v0.3.0

Ordered by how much code it costs.

| Change | Current spec position | This library | Confidence |
|---|---|---|---|
| Card location | `<streamable-http-url>/server-card`; `.well-known` explicitly rejected | `/.well-known/mcp-server-card` | `[HIGH: discovery.md "Hosted Server Card Location" + "Alternatives considered"]` |
| Multi-server | AI Catalog at `/.well-known/ai-catalog.json` with `url` per entry; no path guessing | `WellKnownPathFor()` guesses a sub-path | `[HIGH: discovery.md AI Catalog section]` |
| Media type | `application/mcp-server-card+json` | `application/json` | `[HIGH: discovery.md, schema.ts docstring, go-sdk#1024 constant]` |
| Accept header | client SHOULD send `Accept: application/mcp-server-card+json` | not handled | `[HIGH: discovery.md, normative SHOULD]` |
| Catalog identifiers | `urn:air:{publisher}:{namespace}:{name}` | not modelled | `[HIGH: discovery.md; adopted via PR #31, ADR 0015]` |
| ETag / 304 | recommended as SHOULD; issue #33 open, PR #46 open | not implemented | `[HIGH: issues #33 and #46 open per gh api; go-sdk#1024 already implements it]` |
| Runtime-consistency requirement | new normative "Consistency with Runtime Behavior" section; cards advisory, not authoritative | not documented | `[HIGH: verbatim section in discovery.md, landed via PR #25]` |
| Registry Server/package types | removed, card-only and remote-only | already card-only | `[HIGH: PR #28 by dsp-ant, merged 18-06-2026]` |
| CORS and caching headers | unchanged (`ACAO: *`, `GET`, `Content-Type`; `Cache-Control: public, max-age=3600`) | matches | `[HIGH: discovery.md Security Considerations vs handler.go]` |
| Card field set | matches: `$schema`, `name`, `version`, `description`, `title?`, `websiteUrl?`, `repository?`, `icons?`, `remotes?`, `_meta?` | matches | `[HIGH: field-by-field read of schema.ts against card.go]` |

Two things the library got right and should keep:

- **The struct shape is still correct.** Every field in `schema.ts` maps to a field in `card.go`, with the same JSON names and the same optionality. The v0.3.0 scope reductions (dropping capabilities, requires, authentication, and the MCP resource endpoint) all held. `[HIGH: full read of both files]`
- **`Repository.id` is the one missing optional field.** `schema.ts` documents it as a forge-owned stable identifier used to detect repository-resurrection attacks. `card.go` has `URL`, `Source`, `Subfolder` but no `ID`. `[HIGH: direct comparison of the two type definitions]`

The `$schema` 404 documented in the README is **still true**. `https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json` returns 404 as of 23-07-2026, and `schema.ts` still mandates that exact URL via a regex pattern. `[HIGH: live curl returning 404, plus the `@pattern` in schema.ts]`

## 6. Watch targets

The bead's instruction to switch from the SEP PR to `experimental-ext-server-card` was right. Refine it: the repo is the venue, but three specific things inside it carry the code-change risk.

**Primary, check first:**

1. **`modelcontextprotocol/go-sdk#1024`** (draft, open). The official Go implementation. Merging changes what this library is for. Watch for: ready-for-review flip, merge, and whether the API stays under `mcp/experimental/`. `[HIGH: PR exists, is the declared Go reference implementation]`
2. **`experimental-ext-server-card/docs/discovery.md`** at HEAD. The discovery mechanics, and the file that changed the path out from under v0.3.0. Diff it rather than reading issue threads. `[HIGH: it carries the normative placement and media-type text]`
3. **`experimental-ext-server-card/schema.ts`** at HEAD. Single source of truth for the card shape, stated as such in its own header comment. `[HIGH: verbatim in the file]`

**Secondary, open issues that would each cost code:**

4. **`experimental-ext-server-card#30`** (open, 17-06-2026, vyshnavigadamsetti): optional tool metadata in the card for offline discovery. Would reverse the primitives exclusion. `[HIGH: open per gh api]`
5. **`experimental-ext-server-card#33` and PR `#46`** (open, SamMorrowDrums): ETag plus `If-None-Match` as a SHOULD. Handler change. `[HIGH: open per gh api]`
6. **`experimental-ext-server-card#43`** (open, 20-07-2026, tadasant): Link header and HTML `<link>` catalog discovery. `[HIGH: open per gh api]`
7. **`experimental-ext-server-card#44`** (open, 20-07-2026, tadasant): nested catalogs with depth cap and cycle tracking. `[HIGH: open per gh api]`
8. **`experimental-ext-server-card#13`** (open, 30-05-2026, tadasant): comprehensive auth scenarios back into the card. The v0.3.0 removal of `Authentication` is the thing this could reverse. `[HIGH: open per gh api]`

**External dependency, new since the library was written:**

9. **`Agent-Card/ai-catalog`**. The discovery entrypoint format is now an external standard, not an MCP-owned one. 179 stars, pushed 23-07-2026, actively developed. MCP has already adopted two of its ADRs (0014 `type` rename, 0015 `urn:air:` identifiers). `[HIGH: gh api on the repo; ADR adoption per PRs #31 and #32 in the ext repo]`

**People:**

- **@tadasant** (Tadas Antanavicius). Drove the Extensions Track refactor and files most of the scope issues. Highest-signal single account. `[HIGH: author of #15, #2893, #43, #44, and most open issues]`
- **@SamMorrowDrums** (Sam Morrow). Owns the Go SDK PR, the AI Catalog adoption, and the ETag work. `[HIGH: author of #1024, #42, #33, #46]`
- **@dsp-ant** (David Soria Parra). SEP author and assigned maintainer. Three unanswered bot inactivity nags (18-05, 01-06, 22-06, 13-07-2026, the last at 20 days). Still authored the card-only PR #28 on 11-06-2026, so not absent, just not answering in-thread. `[HIGH: bot comments and PR #28 authorship both from gh api]`

**Retire as a watch target:** `wong2/mcp-server-detector`. Stale since 07-05-2026, 2 stars, targets the obsolete path. `[HIGH: gh api]`

## 7. What this means for the repo

Not a task list, an accurate statement of position. No code was changed by this check.

- The card **struct** is current. The **transport and discovery layer** is a spec generation behind.
- `WellKnownPath`, `WellKnownPathFor`, and `Content-Type: application/json` are the three concrete divergences from normative text, all in `handler.go`.
- The "reference implementation for Go" framing in the README has a live challenger in the official SDK.
- The deferred item (proposing `.well-known/mcp-server-card` to The Website Specification) should now be **dropped rather than deferred**. The path was considered and rejected by name. Proposing it would propose something the spec's own "Alternatives considered" section argues against. `[HIGH: discovery.md "Alternatives considered"]`

## Confidence distribution

Counted over the 54 tagged claims in this document (verified by grep, not estimated).

| Tier | Count | Share |
|---|---|---|
| HIGH | 47 | 87% |
| MEDIUM | 6 | 11% |
| LOW | 1 | 2% |

Mostly HIGH, as expected for a check run against primary sources: GitHub API responses for state, raw file fetches for normative text, and live HTTP for endpoint behavior. Most claims are cross-referenced by construction (an issue's state plus the file change it produced).

The MEDIUM claims are the single-implementer reports from b-gutman, and the inferences about where the WG routed the value-metadata question. The one LOW claim is "no second production adopter found", which rests on a bounded search rather than a systematic one. It should be read as absence of evidence.

## Method

- GitHub state, issues, PRs, comments, and file lists via `gh api` and `gh pr view` (per `rules/api-querying.md`; no WebFetch was used for JSON).
- Normative spec text fetched raw via `gh api ... -H "Accept: application/vnd.github.raw"` for `docs/discovery.md` and `schema.ts` at HEAD.
- Endpoint behavior via `curl` with status and redirect capture.
- Comparison against `servercard/card.go`, `servercard/handler.go`, `servercard/options.go` at the working-tree state on 23-07-2026.
