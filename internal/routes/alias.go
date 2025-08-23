package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type aliasSetter interface {
	SetAlias(alias string, code string) error
}

type aliasDeleter interface {
	DeleteAlias(alias string) error
}

type AliasSetRequest struct {
	Alias      string `path:"alias" minLength:"1" description:"Alias to set to the dictionary"`
	Dictionary string `json:"dictionary"`
}

func aliasSet(registry aliasSetter) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input AliasSetRequest, output *EmptyResponse) error {
		err := registry.SetAlias(input.Alias, input.Dictionary)
		if errors.Is(spellchecker.ErrNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		return nil
	})

	u.SetTitle("Set dictionary alias")
	u.SetDescription("Assigns an alias to a dictionary. If the alias is already used by another dictionary, it will be reassigned to the current one. This route can be used, for example, to manage dictionary versioning.")
	u.SetExpectedErrors(status.Internal, status.NotFound)

	return u
}

type AliasDeleteRequest struct {
	Alias string `path:"alias" minLength:"1" description:"Alias to delete"`
}

func aliasDelete(registry aliasDeleter) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input AliasDeleteRequest, output *EmptyResponse) error {
		err := registry.DeleteAlias(input.Alias)
		if errors.Is(spellchecker.ErrNotFound, err) {
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
