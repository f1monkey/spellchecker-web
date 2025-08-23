package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type dictionarySaver interface {
	Save(code string) error
}

type DictionarySaveRequest struct {
	Code string `path:"code" minLength:"1"`
}

func dictionarySave(registry dictionarySaver) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DictionarySaveRequest, output *Empty) error {
		err := registry.Save(input.Code)
		if errors.Is(spellchecker.ErrNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		return nil
	})

	u.SetTitle("Save a dictionary")
	u.SetDescription("Forces saving the specified dictionary to disk by its code")
	u.SetExpectedErrors(status.Internal, status.NotFound, status.InvalidArgument)

	return u
}
