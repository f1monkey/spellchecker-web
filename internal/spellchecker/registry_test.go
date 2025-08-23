package spellchecker

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/f1monkey/spellchecker"
	"github.com/stretchr/testify/require"
)

func Test_NewRegistry(t *testing.T) {
	t.Parallel()

	createTestFile := func(t *testing.T, dir string, name string) {
		t.Helper()

		sc, err := spellchecker.New("abc")
		require.NoError(t, err)

		item := RegistryItem{
			Spellchecker: sc,
		}

		data, err := json.Marshal(&item)
		require.NoError(t, err)

		err = os.WriteFile(path.Join(dir, name+extension), data, 0755)
		require.NoError(t, err)
	}

	t.Run("no files", func(t *testing.T) {
		dir := t.TempDir()

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Empty(t, result.items)
	})

	t.Run("one file", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "code")

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Contains(t, result.items, "code")
	})

	t.Run("two files, ok", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "code1")
		createTestFile(t, dir, "code2")

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Contains(t, result.items, "code1")
		require.Contains(t, result.items, "code2")
	})

	t.Run("two files, one has invalid extension", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "code1")

		f, err := os.Create(path.Join(dir, "code2.txt"))
		require.NoError(t, err)
		f.Close()

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Contains(t, result.items, "code1")
		require.NotContains(t, result.items, "code2")
	})

	t.Run("two files, one is corrupted", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "code1")

		err := os.WriteFile(path.Join(dir, "code2"+extension), []byte("qweqwe"), 0755)
		require.NoError(t, err)

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Contains(t, result.items, "code1")
		require.NotContains(t, result.items, "code2")
	})
}

func Test_Registry_Add(t *testing.T) {
	t.Parallel()

	t.Run("already exists", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		r.items["code"] = RegistryItem{}

		_, err = r.Add("code", Options{Alphabet: "abc"})

		require.ErrorIs(t, err, ErrAlreadyExists)
	})

	t.Run("spellchecker init error", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "aaa"})

		require.ErrorIs(t, err, ErrSpellcheckerInit)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		result, err := r.Add("code", Options{Alphabet: "abc"})

		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func Test_Registry_Get(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		require.NotNil(t, r.Get("code"))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		require.Nil(t, r.Get("code"))
	})
}
