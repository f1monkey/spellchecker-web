package spellchecker

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Registry_Save(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.Save("qwerty")
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

		err = r.Save(code)
		require.NoError(t, err)
		require.FileExists(t, path.Join(dir, fileName(code)))

		r2, err := NewRegistry(context.Background(), dir)
		require.NoError(t, err)
		require.Contains(t, r2.items, code)
	})
}

func Test_Registry_SaveAll(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		code := "code"

		r, err := NewRegistry(context.Background(), dir)
		require.NoError(t, err)

		_, err = r.Add(code, Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.SaveAll(context.Background())
		require.NoError(t, err)
		require.FileExists(t, path.Join(dir, fileName(code)))
		require.FileExists(t, path.Join(dir, metadataFile))

		r2, err := NewRegistry(context.Background(), dir)
		require.NoError(t, err)
		require.Contains(t, r2.items, code)
	})
}
