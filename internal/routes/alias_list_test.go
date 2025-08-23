package routes

import (
	"context"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
)

type testAliasLister struct {
	items []spellchecker.ListItem
}

func (f *testAliasLister) ListAliases() []spellchecker.ListItem {
	return f.items
}

func Test_AliasList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		lister    *testAliasLister
		wantItems []ListItem
	}{
		{
			name:      "empty list",
			lister:    &testAliasLister{items: []spellchecker.ListItem{}},
			wantItems: []ListItem{},
		},
		{
			name: "single item",
			lister: &testAliasLister{items: []spellchecker.ListItem{
				{Code: "en", Aliases: []string{"eng", "english"}},
			}},
			wantItems: []ListItem{
				{Code: "en", Aliases: []string{"eng", "english"}},
			},
		},
		{
			name: "multiple items",
			lister: &testAliasLister{items: []spellchecker.ListItem{
				{Code: "en", Aliases: []string{"eng"}},
				{Code: "fr", Aliases: []string{"fra", "french"}},
			}},
			wantItems: []ListItem{
				{Code: "en", Aliases: []string{"eng"}},
				{Code: "fr", Aliases: []string{"fra", "french"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := aliasList(tt.lister)

			var out AliasListResponse
			err := interactor.Interact(context.Background(), Empty{}, &out)

			require.NoError(t, err)
			require.Equal(t, tt.wantItems, out.Items)
		})
	}
}
