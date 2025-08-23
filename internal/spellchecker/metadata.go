package spellchecker

import "slices"

const metadataFile = "metadata"

type Metadata struct {
	Aliases         map[string]string   `json:"aliases"`         // alias => dict
	InvertedAliases map[string][]string `json:"invertedAliases"` // dict => aliases
}

func newMetadata() Metadata {
	return Metadata{
		Aliases:         make(map[string]string),
		InvertedAliases: make(map[string][]string),
	}
}

func (r *Registry) SetAlias(alias string, to string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[to]
	if !ok {
		return ErrNotFound
	}

	existing, ok := r.metadata.Aliases[alias]
	if ok {
		if existing == to {
			return nil
		}

		for i, a := range r.metadata.InvertedAliases[existing] {
			if a != alias {
				continue
			}

			r.metadata.InvertedAliases[existing] = slices.Delete(r.metadata.InvertedAliases[existing], i, i+1)
		}
	}

	r.metadata.Aliases[alias] = to
	r.metadata.InvertedAliases[to] = append(r.metadata.InvertedAliases[to], alias)

	if err := r.doSaveMetadata(); err != nil {
		return err
	}

	return nil
}
