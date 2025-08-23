package routes

import (
	"context"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type dictionaryLister interface {
	List() []spellchecker.ListItem
}

type DictionaryListResponse struct {
	Items []ListItem `json:"items"`
}

type ListItem struct {
	Code    string   `json:"code"`
	Aliases []string `json:"aliases"`
}

func dictionaryList(registry dictionaryLister) usecase.Interactor {
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
