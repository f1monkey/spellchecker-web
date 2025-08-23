package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type DictionaryListResponse struct {
	Items []ListItem `json:"items"`
}

type ListItem struct {
	Code    string   `json:"code"`
	Aliases []string `json:"aliases"`
}

func dictionaryList(registry *spellchecker.Registry) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input Empty, output *DictionaryListResponse) error {
		items := registry.List()

		result := make([]ListItem, 0, len(items))

		for _, item := range items {
			result = append(result, ListItem{
				Code:    item.Code,
				Aliases: item.Aliases,
			})
		}

		output.Items = result

		return nil
	})

	u.SetTitle("List all dictionaries")
	u.SetDescription("With their aliases")
	u.SetExpectedErrors(status.Internal, status.AlreadyExists, status.InvalidArgument)

	return u
}

type DictionaryCreateRequest struct {
	Code string `path:"code" minLength:"1"`

	Alphabet  string `json:"alphabet" minLength:"1"`
	MaxErrors uint   `json:"maxErrors" minimum:"0" maximum:"5"`
}

func dictionaryCreate(registry *spellchecker.Registry) usecase.Interactor {
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

type DictionaryDeleteRequest struct {
	Code string `path:"code" minLength:"1"`
}

func dictionaryDelete(registry *spellchecker.Registry) usecase.Interactor {
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

type DictionarySaveRequest struct {
	Code string `path:"code" minLength:"1"`
}

func dictionarySave(registry *spellchecker.Registry) usecase.Interactor {
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
