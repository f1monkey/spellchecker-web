package spellchecker

import (
	"bytes"
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
	ErrAlreadyExists    = fmt.Errorf("already exists")
	ErrSpellcheckerInit = fmt.Errorf("spellchecker init err")
)

const extension = ".dict"

type Registry struct {
	mu sync.RWMutex

	dir   string
	items map[string]*spellchecker.Spellchecker
}

func NewRegistry(ctx context.Context, dir string) (*Registry, error) {
	files, err := listFiles(dir)
	if err != nil {
		return nil, err
	}

	items := make(map[string]*spellchecker.Spellchecker)

	for _, f := range files {
		buf, err := os.ReadFile(path.Join(dir, f.Name()))
		if err != nil {
			logger.FromContext(ctx).Error("registry: read file err", "file", f, "error", err)
		}

		sc, err := spellchecker.Load(bytes.NewBuffer(buf))
		if err != nil {
			logger.FromContext(ctx).Error("registry: unable to initialize spellchecker", "file", f, "error", err)
			continue
		}

		code, _ := strings.CutSuffix(f.Name(), extension)

		items[code] = sc
	}

	return &Registry{
		items: items,
	}, nil
}

func (r *Registry) Add(code string, alphabet string) (*spellchecker.Spellchecker, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[code]; ok {
		return nil, ErrAlreadyExists
	}

	result, err := spellchecker.New(alphabet)
	if err != nil {
		return nil, ErrSpellcheckerInit
	}

	r.items[code] = result

	return result, nil
}

func (r *Registry) Get(code string) *spellchecker.Spellchecker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.items[code]
}

func listFiles(dir string) ([]fs.DirEntry, error) {
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
