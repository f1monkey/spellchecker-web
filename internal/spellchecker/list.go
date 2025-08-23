package spellchecker

type ListItem struct {
	Code    string
	Aliases []string
}

func (r *Registry) List() []ListItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ListItem, 0, len(r.items))

	for code := range r.items {
		result = append(result, ListItem{
			Code:    code,
			Aliases: r.metadata.InvertedAliases[code],
		})
	}

	return result
}
