package routes

import (
	"context"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
)

// тестовый dictionaryLister
type testDictionaryLister struct {
	items []spellchecker.ListItem
}

func (f *testDictionaryLister) List() []spellchecker.ListItem {
	return f.items
}

func Test_DictionaryList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		lister    *testDictionaryLister
		wantItems []ListItem
	}{
		{
			name:      "empty list",
			lister:    &testDictionaryLister{items: []spellchecker.ListItem{}},
			wantItems: []ListItem{},
		},
		{
			name: "single item",
			lister: &testDictionaryLister{items: []spellchecker.ListItem{
				{Code: "en", Aliases: []string{"eng", "english"}},
			}},
			wantItems: []ListItem{
				{Code: "en", Aliases: []string{"eng", "english"}},
			},
		},
		{
			name: "multiple items",
			lister: &testDictionaryLister{items: []spellchecker.ListItem{
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := dictionaryList(tt.lister)

			var out DictionaryListResponse
			err := interactor.Interact(context.Background(), Empty{}, &out)

			require.NoError(t, err)
			require.Equal(t, tt.wantItems, out.Items)
		})
	}
}
