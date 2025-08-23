package routes

import (
	"context"

	"github.com/f1monkey/spellchecker"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type DictionaryAddRequest struct {
	Phrases []DictionaryPhrase `json:"phrases"`
}

type DictionaryPhrase struct {
	Text   string `json:"text" description:"The word or phrase to be added to the dictionary."`
	Weight uint   `json:"weight" min:"1" description:"A numeric value indicating the importance or influence of this entry in spellchecking or suggestions."`
}

type DictionaryAddResponse struct {
	Words int `json:"words" description:"Number of phrases successfully added."`
}

func DictionaryAdd(sc *spellchecker.Spellchecker) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DictionaryAddRequest, output *DictionaryAddResponse) error {
		if len(input.Phrases) == 0 {
			return nil
		}

		wordCnt := 0

		for i := range input.Phrases {

			words := wordSymbols.FindAllString(input.Phrases[i].Text, -1)
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
