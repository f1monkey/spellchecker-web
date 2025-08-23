package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type aliasDeleter interface {
	DeleteAlias(alias string) error
}

type AliasDeleteRequest struct {
	Alias string `path:"alias" minLength:"1" description:"Alias to delete"`
}

func aliasDelete(registry aliasDeleter) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input AliasDeleteRequest, output *Empty) error {
		err := registry.DeleteAlias(input.Alias)
		if errors.Is(spellchecker.ErrAliasNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		return nil
	})

	u.SetTitle("Delete alias from a dictionary")
	u.SetDescription("Removes an alias from a dictionary.")
	u.SetExpectedErrors(status.Internal, status.NotFound)

	return u
}
