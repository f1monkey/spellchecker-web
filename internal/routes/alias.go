package routes

import (
	"context"
	"errors"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type aliasLister interface {
	ListAliases() []spellchecker.ListItem
}

type aliasGetter interface {
	GetCodeByAlias(alias string) (string, error)
}

type aliasSetter interface {
	SetAlias(alias string, code string) error
}

type aliasDeleter interface {
	DeleteAlias(alias string) error
}

type AliasListResponse struct {
	Items []ListItem `json:"items"`
}

func aliasList(registry aliasLister) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input Empty, output *AliasListResponse) error {
		items := registry.ListAliases()

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

	u.SetTitle("List all aliases")
	u.SetDescription("With their aliasesdictionaries")
	u.SetExpectedErrors(status.Internal, status.AlreadyExists, status.InvalidArgument)

	return u
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

type AliasSetRequest struct {
	Alias      string `path:"alias" minLength:"1" description:"Alias to set to the dictionary"`
	Dictionary string `json:"dictionary"`
}

func aliasSet(registry aliasSetter) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input AliasSetRequest, output *Empty) error {
		err := registry.SetAlias(input.Alias, input.Dictionary)
		if errors.Is(spellchecker.ErrAliasNotFound, err) {
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
