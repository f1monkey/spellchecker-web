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

	t.Run("no files", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Empty(t, result.items)
	})

	t.Run("one file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		createTestFile(t, dir, "code")

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Contains(t, result.items, "code")
	})

	t.Run("two files, has metadata", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		createTestFile(t, dir, "code1")
		createTestFile(t, dir, "code2")

		err := os.WriteFile(path.Join(dir, metadataFile), []byte(`{}`), 0644)
		require.NoError(t, err)

		result, err := NewRegistry(context.Background(), dir)

		require.NoError(t, err)
		require.Contains(t, result.items, "code1")
		require.Contains(t, result.items, "code2")
	})

	t.Run("two files, one has invalid extension", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		dir := t.TempDir()

		createTestFile(t, dir, "code1")

		err := os.WriteFile(path.Join(dir, fileName("code2")), []byte("qweqwe"), 0755)
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

		v, err := r.Get("code")
		require.NoError(t, err)

		require.NotNil(t, v)
	})

	t.Run("by alias", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		r.SetAlias("code2", "code")

		v, err := r.Get("code2")
		require.NoError(t, err)
		require.NotNil(t, v)

		v, err = r.Get("code")
		require.NoError(t, err)
		require.NotNil(t, v)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		v, err := r.Get("code")
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, v)
	})
}

func Test_Registry_Delete(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.Delete("qwerty")
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		code := "code"

		r, err := NewRegistry(context.Background(), dir)
		require.NoError(t, err)

		_, err = r.Add(code, Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.Delete(code)
		require.NoError(t, err)

		require.NotContains(t, r.items, code)
	})

	t.Run("success, delete file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		code := "code"

		r, err := NewRegistry(context.Background(), dir)
		require.NoError(t, err)

		_, err = r.Add(code, Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.Save(code)
		require.NoError(t, err)
		require.FileExists(t, path.Join(dir, fileName(code)))

		err = r.Delete(code)
		require.NoError(t, err)
		require.NoFileExists(t, path.Join(dir, fileName(code)))

		require.NotContains(t, r.items, code)
	})
}

func createTestFile(t *testing.T, dir string, name string) {
	t.Helper()

	sc, err := spellchecker.New("abc")
	require.NoError(t, err)

	item := RegistryItem{
		Spellchecker: sc,
	}

	data, err := json.Marshal(&item)
	require.NoError(t, err)

	err = os.WriteFile(path.Join(dir, fileName(name)), data, 0755)
	require.NoError(t, err)
}
