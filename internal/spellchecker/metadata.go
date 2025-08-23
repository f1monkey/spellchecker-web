package spellchecker

import (
	"fmt"
	"slices"
)

const metadataFile = "metadata"

var (
	ErrAliasNotFound = fmt.Errorf("alias not found")
)

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

func (r *Registry) ListAliases() []ListItem {
	return r.List()
}

func (r *Registry) GetCodeByAlias(alias string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.metadata.Aliases[alias]
	if !ok {
		return "", ErrAliasNotFound
	}

	return v, nil
}

func (r *Registry) SetAlias(alias string, to string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[to]
	if !ok {
		return ErrAliasNotFound
	}

	existing, ok := r.metadata.Aliases[alias]
	if ok {
		if existing == to {
			return nil
		}

		_ = r.doDeleteAlias(alias)
	}

	r.metadata.Aliases[alias] = to
	r.metadata.InvertedAliases[to] = append(r.metadata.InvertedAliases[to], alias)

	if err := r.doSaveMetadata(); err != nil {
		return err
	}

	return nil
}

func (r *Registry) DeleteAlias(alias string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.doDeleteAlias(alias)
	if err != nil {
		return err
	}

	if err := r.doSaveMetadata(); err != nil {
		return err
	}

	return nil
}

func (r *Registry) doDeleteAlias(alias string) error {
	existing, ok := r.metadata.Aliases[alias]
	if !ok {
		return ErrAliasNotFound
	}

	delete(r.metadata.Aliases, alias)

	for i, a := range r.metadata.InvertedAliases[existing] {
		if a != alias {
			continue
		}

		r.metadata.InvertedAliases[existing] = slices.Delete(r.metadata.InvertedAliases[existing], i, i+1)

		if len(r.metadata.InvertedAliases[existing]) == 0 {
			delete(r.metadata.InvertedAliases, existing)
		}
	}

	return nil
}
