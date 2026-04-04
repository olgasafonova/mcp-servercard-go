package servercard

import (
	"errors"
	"strings"
)

// Build constructs a ServerCard from the provided Options.
// It validates that required fields (Name, Version, Description) are set
// and that Name follows the reverse-DNS format with exactly one forward slash.
func Build(opts Options) (*ServerCard, error) {
	if opts.Name == "" {
		return nil, errors.New("servercard: name is required")
	}
	if strings.Count(opts.Name, "/") != 1 {
		return nil, errors.New("servercard: name must contain exactly one forward slash (reverse-DNS format)")
	}
	if opts.Version == "" {
		return nil, errors.New("servercard: version is required")
	}
	if opts.Description == "" {
		return nil, errors.New("servercard: description is required")
	}

	card := &ServerCard{
		Schema:       SchemaURL,
		Name:         opts.Name,
		Version:      opts.Version,
		Description:  opts.Description,
		Title:        opts.Title,
		WebsiteURL:   opts.WebsiteURL,
		Repository:   opts.Repository,
		Icons:        opts.Icons,
		Remotes:      opts.Remotes,
		Capabilities: opts.Capabilities,
		Requires:     opts.Requires,
		Meta:         opts.Meta,
	}

	// Normalize Auth.Schemes to avoid null in JSON output.
	for i := range card.Remotes {
		if card.Remotes[i].Authentication != nil {
			card.Remotes[i].Authentication.Normalize()
		}
	}

	// Store provider info in _meta if provided.
	if opts.Provider != nil {
		if card.Meta == nil {
			card.Meta = make(map[string]any)
		}
		card.Meta["provider"] = opts.Provider
	}

	return card, nil
}
