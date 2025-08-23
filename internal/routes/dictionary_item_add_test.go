package routes

import (
	"context"
	"errors"
	"regexp"
	"testing"

	f1mspellchecker "github.com/f1monkey/spellchecker"
	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testDictionary struct {
	added [][]string
	calls []uint
}

func (d *testDictionary) AddWeight(weight uint, words ...string) {
	d.added = append(d.added, words)
	d.calls = append(d.calls, weight)
}

func Test_DictionaryItemAdd(t *testing.T) {
	t.Parallel()

	splitter := regexp.MustCompile(`[a-zA-Z]+`)

	tests := []struct {
		name       string
		getter     *testDictionaryGetter
		input      DictionaryItemAddRequest
		wantErr    bool
		wantCode   status.Code
		wantWords  int
		wantAdded  [][]string
		wantWeight []uint
	}{
		{
			name: "success single phrase with weight",
			getter: &testDictionaryGetter{
				err: nil,
			},
			input: DictionaryItemAddRequest{
				Code: "en",
				Phrases: []DictionaryItemPhrase{
					{Text: "hello world", Weight: 2},
				},
			},
			wantErr:    false,
			wantCode:   status.OK,
			wantWords:  2,
			wantAdded:  [][]string{{"hello", "world"}},
			wantWeight: []uint{2},
		},
		{
			name: "phrase with zero weight gets default=1",
			getter: &testDictionaryGetter{
				err: nil,
			},
			input: DictionaryItemAddRequest{
				Code: "en",
				Phrases: []DictionaryItemPhrase{
					{Text: "hi", Weight: 0},
				},
			},
			wantErr:    false,
			wantCode:   status.OK,
			wantWords:  1,
			wantAdded:  [][]string{{"hi"}},
			wantWeight: []uint{1},
		},
		{
			name: "phrase with no words (ignored)",
			getter: &testDictionaryGetter{
				err: nil,
			},
			input: DictionaryItemAddRequest{
				Code: "en",
				Phrases: []DictionaryItemPhrase{
					{Text: "!!!", Weight: 5}, // regex не найдёт слов
				},
			},
			wantErr:    false,
			wantCode:   status.OK,
			wantWords:  0,
			wantAdded:  nil,
			wantWeight: nil,
		},
		{
			name: "dictionary not found",
			getter: &testDictionaryGetter{
				err: spellchecker.ErrNotFound,
			},
			input:    DictionaryItemAddRequest{Code: "xx"},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name: "internal error",
			getter: &testDictionaryGetter{
				err: errors.New("boom"),
			},
			input:    DictionaryItemAddRequest{Code: "en"},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc, err := f1mspellchecker.New(f1mspellchecker.DefaultAlphabet)
			require.NoError(t, err)

			tt.getter.sc = sc

			interactor := dictionaryItemAdd(tt.getter, splitter)

			var out DictionaryItemAddResponse
			err = interactor.Interact(context.Background(), tt.input, &out)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, err.(isErr).Is(tt.wantCode))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantWords, out.Words)
			}
		})
	}
}
