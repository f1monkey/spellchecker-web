package routes

import (
	"context"
	"errors"
	"regexp"
	"testing"

	f1mspellchecker "github.com/f1monkey/spellchecker"
	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testDictionaryGetter struct {
	sc  *f1mspellchecker.Spellchecker
	err error
}

func (f *testDictionaryGetter) Get(code string) (*f1mspellchecker.Spellchecker, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.sc, nil
}

func Test_DictionaryFix(t *testing.T) {
	t.Parallel()

	splitter := regexp.MustCompile(`[a-zA-Z]+`)

	sc, err := f1mspellchecker.New(f1mspellchecker.DefaultAlphabet)
	require.NoError(t, err)

	sc.Add("hello")

	tests := []struct {
		name      string
		getter    dictionaryGetter
		input     DictionaryFixRequest
		wantErr   bool
		wantCode  status.Code
		wantFixes []Fix
	}{
		{
			name:      "empty text",
			getter:    &testDictionaryGetter{sc: &f1mspellchecker.Spellchecker{}},
			input:     DictionaryFixRequest{Code: "en", Text: "", Limit: 5},
			wantErr:   false,
			wantFixes: []Fix{},
		},
		{
			name: "exact match word",
			getter: &testDictionaryGetter{
				sc: sc,
			},
			input:     DictionaryFixRequest{Code: "en", Text: "hello", Limit: 5},
			wantErr:   false,
			wantFixes: []Fix{},
		},
		{
			name: "word with suggestions",
			getter: &testDictionaryGetter{
				sc: sc,
			},
			input:   DictionaryFixRequest{Code: "en", Text: "hellp", Limit: 5},
			wantErr: false,
			wantFixes: []Fix{
				{
					Start: 0, End: 5,
					Error: "invalid_word",
					Suggestions: []SpellcheckerSuggestion{
						{Text: "hello"},
					},
				},
			},
		},
		{
			name: "word without suggestions",
			getter: &testDictionaryGetter{
				sc: sc,
			},
			input:   DictionaryFixRequest{Code: "en", Text: "qwertyuiop", Limit: 5},
			wantErr: false,
			wantFixes: []Fix{
				{
					Start: 0, End: 10,
					Error: "unknown_word",
				},
			},
		},
		{
			name:     "dictionary not found",
			getter:   &testDictionaryGetter{err: spellchecker.ErrNotFound},
			input:    DictionaryFixRequest{Code: "xx", Text: "hello", Limit: 5},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name:     "internal error",
			getter:   &testDictionaryGetter{err: errors.New("boom")},
			input:    DictionaryFixRequest{Code: "en", Text: "hello", Limit: 5},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := dictionaryFix(tt.getter, splitter)

			var out DictionaryFixResponse
			err := interactor.Interact(context.Background(), tt.input, &out)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, err.(isErr).Is(tt.wantCode))
			} else {
				require.NoError(t, err)

				require.Len(t, out.Fixes, len(tt.wantFixes))

				for i, f := range tt.wantFixes {

					require.Len(t, out.Fixes[i].Suggestions, len(f.Suggestions))

					for j, s := range f.Suggestions {
						assert.Equal(t, s.Text, tt.wantFixes[i].Suggestions[j].Text)
					}

					assert.Equal(t, f.Start, out.Fixes[i].Start)
					assert.Equal(t, f.End, out.Fixes[i].End)
					assert.Equal(t, f.Error, out.Fixes[i].Error)
				}
			}
		})
	}
}
