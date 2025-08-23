package routes

import (
	"context"
	"errors"

	f1mspellchecker "github.com/f1monkey/spellchecker"
	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type registryAdder interface {
	Add(code string, options spellchecker.Options) (*f1mspellchecker.Spellchecker, error)
}

type DictionaryCreateRequest struct {
	Code string `path:"code" minLength:"1"`

	Alphabet  string `json:"alphabet" minLength:"1"`
	MaxErrors uint   `json:"maxErrors" minimum:"0" maximum:"5"`
}

func dictionaryCreate(registry registryAdder) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DictionaryCreateRequest, output *Empty) error {
		_, err := registry.Add(input.Code, spellchecker.Options{
			Alphabet:  input.Alphabet,
			MaxErrors: input.MaxErrors,
		})
		if errors.Is(spellchecker.ErrAlreadyExists, err) {
			return status.Wrap(err, status.AlreadyExists)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		return nil
	})

	u.SetTitle("Create a new dictionary")
	u.SetDescription("Adds a new dictionary to the registry")
	u.SetExpectedErrors(status.Internal, status.AlreadyExists, status.InvalidArgument)

	return u
}
