package spellchecker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/f1monkey/spellchecker"
	"github.com/f1monkey/spellchecker-web/internal/logger"
)

var (
	ErrAlreadyExists    = fmt.Errorf("dictionary already exists")
	ErrSpellcheckerInit = fmt.Errorf("spellchecker init err")
	ErrNotFound         = fmt.Errorf("dictionary not found")
)

const extension = ".dict"

type Registry struct {
	mu sync.RWMutex

	dir   string
	items map[string]RegistryItem
}

func NewRegistry(ctx context.Context, dir string) (*Registry, error) {
	files, err := findDictionaries(dir)
	if err != nil {
		return nil, err
	}

	items := make(map[string]RegistryItem)

	for _, f := range files {
		buf, err := os.ReadFile(path.Join(dir, f.Name()))
		if err != nil {
			logger.FromContext(ctx).Error("registry: read file err", "file", f, "error", err)
			continue
		}

		var item RegistryItem

		if err := json.Unmarshal(buf, &item); err != nil {
			logger.FromContext(ctx).Error("registry: unable to initalize registry item", "file", f, "error", err)
			continue
		}

		code, _ := strings.CutSuffix(f.Name(), extension)

		logger.FromContext(ctx).Info("registry: loaded dictionary", "code", code)

		items[code] = item
	}

	return &Registry{
		items: items,
	}, nil
}

func (r *Registry) Add(code string, options Options) (*spellchecker.Spellchecker, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[code]; ok {
		return nil, ErrAlreadyExists
	}

	result, err := spellchecker.New(
		options.Alphabet,
		spellchecker.WithMaxErrors(int(options.MaxErrors)),
	)
	if err != nil {
		return nil, ErrSpellcheckerInit
	}

	r.items[code] = RegistryItem{
		Spellchecker: result,
	}

	return result, nil
}

func (r *Registry) Get(code string) *spellchecker.Spellchecker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.items[code].Spellchecker
}

func (r *Registry) Delete(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[code]
	if !ok {
		return ErrNotFound
	}

	delete(r.items, code)

	return nil
}

func findDictionaries(dir string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make([]fs.DirEntry, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), extension) {
			continue
		}

		result = append(result, f)
	}

	return result, nil
}
