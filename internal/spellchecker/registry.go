package spellchecker

import (
	"context"
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

	metadata Metadata
	dir      string
	items    map[string]RegistryItem
}

func NewRegistry(ctx context.Context, dir string) (*Registry, error) {
	files, err := findDictionaries(dir)
	if err != nil {
		return nil, err
	}

	result := &Registry{
		dir:   dir,
		items: make(map[string]RegistryItem),
	}

	metadata, err := result.doLoadMetadata()
	if err != nil {
		return nil, err
	}

	result.metadata = metadata

	for _, f := range files {
		code, _ := strings.CutSuffix(f.Name(), extension)

		item, err := result.doLoad(code)
		if err != nil {
			logger.FromContext(ctx).Error("registry: dictionary load error", "code", code, "error", err)
			continue
		}

		logger.FromContext(ctx).Info("registry: loaded dictionary", "dictionary", code)

		result.items[code] = item
	}

	return result, nil
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
		Options:      options,
	}

	return result, nil
}

func (r *Registry) Get(code string) (*spellchecker.Spellchecker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.items[code]
	if !ok {
		aliased, ok := r.metadata.Aliases[code]
		if !ok {
			return nil, ErrNotFound
		}

		v = r.items[aliased]
	}

	return v.Spellchecker, nil
}

func (r *Registry) Delete(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[code]
	if !ok {
		return ErrNotFound
	}

	err := os.Remove(fullPath(r.dir, code))
	if err != nil && !os.IsNotExist(err) {
		return err
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

func fileName(code string) string {
	return code + extension
}

func fullPath(dir string, code string) string {
	return path.Join(dir, fileName(code))
}
