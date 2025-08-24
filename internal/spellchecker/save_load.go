package spellchecker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/f1monkey/spellchecker-web/internal/logger"
)

func (r *Registry) AutoSave(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := r.SaveAll(ctx); err != nil {
					logger.FromContext(ctx).Error("registry: save all error", "error", err)
				}
			}
		}
	}()
}

func (r *Registry) SaveAll(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := r.doSaveMetadata(); err != nil {
		return fmt.Errorf("metadata save: %w", err)
	}

	for code := range r.items {
		if err := r.doSave(code); err != nil {
			return fmt.Errorf("dictionary %q save: %w", code, err)
		}

		logger.FromContext(ctx).Info("registry: dictionary saved", "dictionary", code)
	}

	return nil
}

func (r *Registry) Save(code string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.doSave(code)
}

func (r *Registry) doSave(code string) error {
	item, ok := r.items[code]
	if !ok {
		return ErrNotFound
	}

	data, err := json.Marshal(&item)
	if err != nil {
		return err
	}

	dstPath := path.Join(r.dir, fileName(code))

	tmpFile, err := os.CreateTemp(r.dir, fileName(code)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	return os.Rename(tmpName, dstPath)
}

func (r *Registry) doLoad(code string) (RegistryItem, error) {
	buf, err := os.ReadFile(fullPath(r.dir, code))
	if err != nil {
		return RegistryItem{}, err
	}

	var item RegistryItem

	if err := json.Unmarshal(buf, &item); err != nil {
		return RegistryItem{}, err
	}

	return item, nil
}

func (r *Registry) doLoadMetadata() (Metadata, error) {
	data, err := os.ReadFile(path.Join(r.dir, metadataFile))
	if os.IsNotExist(err) {
		return newMetadata(), nil
	} else if err != nil {
		return newMetadata(), err
	}

	var result Metadata

	return result, json.Unmarshal(data, &result)
}

func (r *Registry) doSaveMetadata() error {
	data, err := json.Marshal(r.metadata)
	if err != nil {
		return err
	}

	dstPath := path.Join(r.dir, metadataFile)

	tmpFile, err := os.CreateTemp(r.dir, metadataFile+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	return os.Rename(tmpName, dstPath)
}
