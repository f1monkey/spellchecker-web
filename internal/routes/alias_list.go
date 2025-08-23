package routes

import (
	"context"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type aliasLister interface {
	ListAliases() []spellchecker.ListItem
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
