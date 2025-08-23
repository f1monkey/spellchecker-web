package routes

import (
	"context"
	"errors"
	"regexp"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type DictionaryItemAddRequest struct {
	Code string `path:"code" minLength:"1"`

	Phrases []DictionaryItemPhrase `json:"phrases" minLength:"1"`
}

type DictionaryItemPhrase struct {
	Text   string `json:"text" description:"The word or phrase to be added to the dictionary."`
	Weight uint   `json:"weight" min:"1" description:"A numeric value indicating the importance or influence of this entry in spellchecking or suggestions."`
}

type DictionaryItemAddResponse struct {
	Words int `json:"words" description:"Number of phrases successfully added."`
}

func dictionaryItemAdd(registry *spellchecker.Registry, splitter *regexp.Regexp) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DictionaryItemAddRequest, output *DictionaryItemAddResponse) error {
		sc, err := registry.Get(input.Code)
		if errors.Is(spellchecker.ErrNotFound, err) {
			return status.Wrap(err, status.NotFound)
		} else if err != nil {
			return status.Wrap(err, status.Internal)
		}

		wordCnt := 0

		for i := range input.Phrases {

			words := splitter.FindAllString(input.Phrases[i].Text, -1)
			if len(words) == 0 {
				continue
			}

			weight := input.Phrases[i].Weight
			if weight == 0 {
				weight = 1
			}

			sc.AddWeight(weight, words...)
			wordCnt += len(words)
		}

		output.Words = wordCnt

		return nil
	})

	u.SetTitle("Add phrases/words to spellchecker")
	u.SetDescription("Adds one or more custom phrases or words to the spellchecker dictionary. Each phrase can have an optional weight to influence matching or prioritization.")
	u.SetExpectedErrors(status.Internal)

	return u
}
