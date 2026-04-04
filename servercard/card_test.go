package servercard_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/olgasafonova/mcp-servercard-go/servercard"
)

func minimalOpts() servercard.Options {
	return servercard.Options{
		Name:        "io.github.test/minimal-server",
		Version:     "1.0.0",
		Description: "A minimal test server.",
	}
}

func fullOpts() servercard.Options {
	return servercard.Options{
		Name:        "io.github.olgasafonova/gleif-mcp-server",
		Version:     "1.4.0",
		Description: "Access the GLEIF LEI database for company verification, KYC, and ownership research.",
		Title:       "GLEIF MCP Server",
		WebsiteURL:  "https://github.com/olgasafonova/gleif-mcp-server",
		Repository: &servercard.Repository{
			URL:    "https://github.com/olgasafonova/gleif-mcp-server",
			Source: "github",
		},
		Icons: []servercard.Icon{
			{Source: "https://example.com/icon.png", Sizes: []string{"48x48"}, MIMEType: "image/png"},
		},
		Remotes: []servercard.Remote{
			{
				Type:                      "streamable-http",
				URL:                       "https://gleif.example.com/mcp",
				SupportedProtocolVersions: []string{"2025-03-12"},
				Authentication:            &servercard.Auth{Required: false, Schemes: []string{}},
			},
		},
		Capabilities: &servercard.Capabilities{
			Tools:   &servercard.ToolsCap{ListChanged: false},
			Prompts: &servercard.PromptsCap{ListChanged: false},
		},
		Provider: &servercard.Provider{
			Name: "Olga Safonova",
			URL:  "https://github.com/olgasafonova",
		},
	}
}

func TestBuildMinimal(t *testing.T) {
	card := servercard.Build(minimalOpts())

	if card.Schema != servercard.SchemaURL {
		t.Errorf("Schema = %q, want %q", card.Schema, servercard.SchemaURL)
	}
	if card.Name != "io.github.test/minimal-server" {
		t.Errorf("Name = %q", card.Name)
	}
	if card.Version != "1.0.0" {
		t.Errorf("Version = %q", card.Version)
	}
	if card.Title != "" {
		t.Errorf("Title should be empty, got %q", card.Title)
	}
	if card.Capabilities != nil {
		t.Error("Capabilities should be nil for minimal card")
	}
}

func TestBuildFull(t *testing.T) {
	card := servercard.Build(fullOpts())

	if card.Title != "GLEIF MCP Server" {
		t.Errorf("Title = %q", card.Title)
	}
	if len(card.Remotes) != 1 {
		t.Fatalf("Remotes len = %d, want 1", len(card.Remotes))
	}
	if card.Remotes[0].Type != "streamable-http" {
		t.Errorf("Remote type = %q", card.Remotes[0].Type)
	}
	if card.Capabilities == nil || card.Capabilities.Tools == nil {
		t.Fatal("Capabilities.Tools should not be nil")
	}
	if card.Capabilities.Resources != nil {
		t.Error("Capabilities.Resources should be nil")
	}

	// Provider should be in _meta.
	p, ok := card.Meta["provider"]
	if !ok {
		t.Fatal("_meta.provider missing")
	}
	prov, ok := p.(*servercard.Provider)
	if !ok {
		t.Fatalf("_meta.provider type = %T", p)
	}
	if prov.Name != "Olga Safonova" {
		t.Errorf("provider name = %q", prov.Name)
	}
}

func TestJSON(t *testing.T) {
	card := servercard.Build(minimalOpts())
	data, err := card.JSON()
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal("invalid JSON:", err)
	}

	if parsed["$schema"] != servercard.SchemaURL {
		t.Errorf("$schema = %v", parsed["$schema"])
	}
	if parsed["name"] != "io.github.test/minimal-server" {
		t.Errorf("name = %v", parsed["name"])
	}
	// Omitempty fields should be absent.
	if _, ok := parsed["title"]; ok {
		t.Error("title should be omitted for empty string")
	}
	if _, ok := parsed["capabilities"]; ok {
		t.Error("capabilities should be omitted when nil")
	}
}

func TestJSONFull(t *testing.T) {
	card := servercard.Build(fullOpts())
	data, err := card.JSON()
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal("invalid JSON:", err)
	}

	if parsed["title"] != "GLEIF MCP Server" {
		t.Errorf("title = %v", parsed["title"])
	}
	if parsed["websiteUrl"] == nil {
		t.Error("websiteUrl should be present")
	}

	// Check remotes.
	remotes, ok := parsed["remotes"].([]any)
	if !ok || len(remotes) != 1 {
		t.Fatalf("remotes = %v", parsed["remotes"])
	}
	remote := remotes[0].(map[string]any)
	if remote["type"] != "streamable-http" {
		t.Errorf("remote type = %v", remote["type"])
	}
}

func TestHandlerGET(t *testing.T) {
	card := servercard.Build(minimalOpts())
	handler := servercard.Handler(card)

	req := httptest.NewRequest(http.MethodGet, servercard.WellKnownPath, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	if cors := resp.Header.Get("Access-Control-Allow-Origin"); cors != "*" {
		t.Errorf("CORS origin = %q", cors)
	}
	if cc := resp.Header.Get("Cache-Control"); cc != "public, max-age=3600" {
		t.Errorf("Cache-Control = %q", cc)
	}

	body, _ := io.ReadAll(resp.Body)
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatal("body is not valid JSON:", err)
	}
}

func TestHandlerOPTIONS(t *testing.T) {
	card := servercard.Build(minimalOpts())
	handler := servercard.Handler(card)

	req := httptest.NewRequest(http.MethodOptions, servercard.WellKnownPath, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
	if cors := resp.Header.Get("Access-Control-Allow-Methods"); cors != "GET" {
		t.Errorf("Allow-Methods = %q", cors)
	}
}

func TestHandlerPOST(t *testing.T) {
	card := servercard.Build(minimalOpts())
	handler := servercard.Handler(card)

	req := httptest.NewRequest(http.MethodPost, servercard.WellKnownPath, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Result().StatusCode)
	}
}

func TestRegisterResource(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	card := servercard.Build(minimalOpts())
	servercard.RegisterResource(server, card)

	// Verify the resource is accessible by listing resources through a session.
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	go func() { _ = server.Run(ctx, serverTransport) }()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = session.Close() }()

	resources, err := session.ListResources(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, r := range resources.Resources {
		if r.URI == servercard.ResourceURI {
			found = true
			if r.MIMEType != "application/json" {
				t.Errorf("resource MIMEType = %q", r.MIMEType)
			}
			break
		}
	}
	if !found {
		t.Error("server card resource not found in resources list")
	}

	// Read the resource content.
	result, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: servercard.ResourceURI,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Contents) != 1 {
		t.Fatalf("contents len = %d", len(result.Contents))
	}
	content := result.Contents[0]
	if content.MIMEType != "application/json" {
		t.Errorf("content MIMEType = %q", content.MIMEType)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(content.Text), &parsed); err != nil {
		t.Fatal("resource content is not valid JSON:", err)
	}
	if parsed["name"] != "io.github.test/minimal-server" {
		t.Errorf("name = %v", parsed["name"])
	}
}

func TestAttach(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	handler := servercard.Attach(server, fullOpts())

	// HTTP handler should work.
	req := httptest.NewRequest(http.MethodGet, servercard.WellKnownPath, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d", rec.Result().StatusCode)
	}

	body, _ := io.ReadAll(rec.Result().Body)
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["title"] != "GLEIF MCP Server" {
		t.Errorf("title = %v", parsed["title"])
	}
}
