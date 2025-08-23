package routes

import (
	"context"
	"errors"
	"regexp"
	"unicode/utf8"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type DictionaryFixRequest struct {
	Code string `path:"code" minLength:"1"`

	Text  string `json:"text" description:"Phrase to be checked"`
	Limit int    `json:"limit" default:"5" desciption:"Max suggestions per word"`
}

type DictionaryFixResponse struct {
	Fixes []Fix `json:"fixes" description:"List of detected issues."`
}

type Fix struct {
	Start       int                      `json:"start" description:"Starting character index of the incorrect word in the input."`
	End         int                      `json:"end" description:"Ending character index."`
	Suggestions []SpellcheckerSuggestion `json:"suggestions,omitempty" description:"List of correction suggestions."`
	Error       string                   `json:"error" enum:"unknown_word,invalid_word" description:"Type of detected error. unknown_word - no possible corrections found; invalid_word - the word can be corrected using one of the provided suggestions"`
}

type SpellcheckerSuggestion struct {
	Text  string  `json:"text" descrption:"Suggested corrected word."`
	Score float64 `json:"score" description:"Confidence score of the suggestion."`
}

func dictionaryFix(registry *spellchecker.Registry, splitter *regexp.Regexp) usecase.Interactor {
	const (
		errorUnknownWord = "unknown_word"
		errorInvalidWord = "invalid_word"
	)

	u := usecase.NewInteractor(func(ctx context.Context, input DictionaryFixRequest, output *DictionaryFixResponse) error {
		sc, err := registry.Get(input.Code)
		if errors.Is(spellchecker.ErrNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		if input.Text == "" {
			return nil
		}

		matches := splitter.FindAllStringIndex(input.Text, -1)
		fixes := make([]Fix, 0, len(matches))

		for _, match := range matches {
			startByte, endByte := match[0], match[1]
			startRune := utf8.RuneCountInString(input.Text[:startByte])
			endRune := startRune + utf8.RuneCountInString(input.Text[startByte:endByte])

			fix := Fix{
				Start: startRune,
				End:   endRune,
			}

			word := input.Text[startByte:endByte]

			suggestions := sc.SuggestScore(word, input.Limit)

			if suggestions.ExactMatch {
				continue
			}

			if len(suggestions.Suggestions) == 0 {
				fix.Error = errorUnknownWord
			} else {
				fix.Error = errorInvalidWord
				fix.Suggestions = make([]SpellcheckerSuggestion, 0, len(suggestions.Suggestions))

				for _, s := range suggestions.Suggestions {
					fix.Suggestions = append(fix.Suggestions, SpellcheckerSuggestion{
						Text:  s.Value,
						Score: s.Score,
					})
				}
			}

			fixes = append(fixes, fix)
		}

		output.Fixes = fixes

		return nil
	})

	u.SetTitle("Fix text")
	u.SetDescription("Performs spellchecking on the given input text. Returns misspelled words along with suggested corrections, up to the specified limit per word.")
	u.SetExpectedErrors(status.Internal, status.NotFound)

	return u
}
