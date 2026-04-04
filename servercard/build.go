package servercard

// Build constructs a ServerCard from the provided Options.
func Build(opts Options) *ServerCard {
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

	// Store provider info in _meta if provided.
	if opts.Provider != nil {
		if card.Meta == nil {
			card.Meta = make(map[string]any)
		}
		card.Meta["provider"] = opts.Provider
	}

	return card
}
