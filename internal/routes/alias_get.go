package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type aliasGetter interface {
	GetCodeByAlias(alias string) (string, error)
}

type AliasGetRequest struct {
	Alias string `path:"alias" minLength:"1" description:"Alias to set to the dictionary"`
}

type AliasGetResponse struct {
	Dictionary string `json:"dictionary"`
}

func aliasGet(registry aliasGetter) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input AliasGetRequest, output *AliasGetResponse) error {
		code, err := registry.GetCodeByAlias(input.Alias)
		if errors.Is(spellchecker.ErrNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		output.Dictionary = code

		return nil
	})

	u.SetTitle("Get dictionary alias")
	u.SetDescription("Returns dictionary code assigned to the provided alias")
	u.SetExpectedErrors(status.Internal, status.NotFound)

	return u
}
