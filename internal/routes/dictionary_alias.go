package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type DictionaryAliasRequest struct {
	Code string `path:"code" minLength:"1"`

	Alias string `json:"alias" description:"Alias to set to the dictionary"`
}

func dictionaryAlias(registry *spellchecker.Registry) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DictionaryAliasRequest, output *EmptyResponse) error {
		err := registry.SetAlias(input.Alias, input.Code)
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
