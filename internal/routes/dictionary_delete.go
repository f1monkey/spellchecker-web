package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type dictionaryDeleter interface {
	Delete(code string) error
}

type DictionaryDeleteRequest struct {
	Code string `path:"code" minLength:"1"`
}

func dictionaryDelete(registry dictionaryDeleter) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DictionaryDeleteRequest, output *Empty) error {
		err := registry.Delete(input.Code)
		if errors.Is(spellchecker.ErrNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		return nil
	})

	u.SetTitle("Delete a dictionary")
	u.SetDescription("Removes a dictionary from the registry")
	u.SetExpectedErrors(status.Internal, status.NotFound, status.InvalidArgument)

	return u
}
