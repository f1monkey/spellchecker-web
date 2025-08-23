package spellchecker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Registry_SetAlias(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		err = r.SetAlias("code2", "code")
		require.ErrorIs(t, err, ErrAliasNotFound)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.SetAlias("code2", "code")
		require.NoError(t, err)
	})

	t.Run("multiple aliases", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.SetAlias("alias1", "code")
		require.NoError(t, err)
		require.Contains(t, r.metadata.Aliases, "alias1")

		err = r.SetAlias("alias2", "code")
		require.NoError(t, err)
		require.Contains(t, r.metadata.Aliases, "alias2")
	})

	t.Run("overwrite alias", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code1", Options{Alphabet: "abc"})
		require.NoError(t, err)

		_, err = r.Add("code2", Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.SetAlias("alias", "code1")
		require.NoError(t, err)
		require.Contains(t, r.metadata.Aliases, "alias")
		require.Contains(t, r.metadata.InvertedAliases["code1"], "alias")

		err = r.SetAlias("alias", "code2")
		require.NoError(t, err)
		require.Contains(t, r.metadata.Aliases, "alias")
		require.NotContains(t, r.metadata.InvertedAliases["code1"], "alias")
		require.Contains(t, r.metadata.InvertedAliases["code2"], "alias")
	})
}

func Test_Registry_DeleteAlias(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		err = r.DeleteAlias("code")
		require.ErrorIs(t, err, ErrAliasNotFound)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		r, err := NewRegistry(context.Background(), t.TempDir())
		require.NoError(t, err)

		_, err = r.Add("code", Options{Alphabet: "abc"})
		require.NoError(t, err)

		err = r.SetAlias("alias", "code")
		require.NoError(t, err)

		err = r.DeleteAlias("alias")
		require.NoError(t, err)
		require.NotContains(t, r.metadata.Aliases, "alias")
		require.NotContains(t, r.metadata.InvertedAliases["code"], "alias")
	})
}
