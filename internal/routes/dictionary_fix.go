package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	f1mspellchecker "github.com/f1monkey/spellchecker/v2"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type dictionaryGetter interface {
	Get(code string) (*f1mspellchecker.Spellchecker, error)
}

type DictionaryFixRequest struct {
	Code                string         `path:"code" minLength:"1" description:"Dictionary code to use for spellchecking."`
	Text                string         `json:"text" description:"Input text to be checked and corrected."`
	Limit               int            `json:"limit" default:"5" description:"Maximum number of suggestions to return per word."`
	MaxErrors           int            `json:"maxErrors" default:"2" description:"Maximum number of bit-level differences allowed between the input word and a dictionary word. Examples: deletion=1 bit (proble→problem), insertion=1 bit (problemm→problem), substitution=2 bits (problam→problem), transposition=0 bits (problme→problem). Not recommended to set higher than 2, as it can impact performance."`
	Fuzziness           FuzzinessValue `json:"fuzziness" description:"Maximum allowed Levenshtein edit distance. Allowed values: '0','1','2'... (fixed distance), 'AUTO' (auto by word length, default AUTO:3,6), 'AUTO:low,high' (custom range). See: https://www.elastic.co/docs/reference/elasticsearch/rest-apis/common-options#fuzziness"`
	SimilarityThreshold float64        `json:"similarityThreshold" minimum:"0" maximum:"1" description:"Required similarity ratio between input word and candidate suggestion (0.0–1.0). Example: 0.6 = candidate must be at least 60% similar to input."`
}

type DictionaryFixResponse struct {
	Fixes   []Fix     `json:"fixes" description:"List of detected issues."`
	Correct []Correct `json:"correct" description:"List of correct words."`
}

type Fix struct {
	Start       int                      `json:"start" description:"Starting character index of the incorrect word in the input."`
	End         int                      `json:"end" description:"Ending character index."`
	Suggestions []SpellcheckerSuggestion `json:"suggestions,omitempty" description:"List of correction suggestions."`
	Error       string                   `json:"error" enum:"unknown_word,invalid_word" description:"Type of detected error. unknown_word - no possible corrections found; invalid_word - the word can be corrected using one of the provided suggestions"`
}

type Correct struct {
	Start int `json:"start" description:"Starting character index of the word in the input."`
	End   int `json:"end" description:"Ending character index."`
}

type SpellcheckerSuggestion struct {
	Text  string  `json:"text" descrption:"Suggested corrected word."`
	Score float64 `json:"score" description:"Confidence score of the suggestion."`
}

func dictionaryFix(registry dictionaryGetter, splitter *regexp.Regexp) usecase.Interactor {
	const (
		errorUnknownWord = "unknown_word"
		errorInvalidWord = "invalid_word"
	)

	u := usecase.NewInteractor(
		func(ctx context.Context, input DictionaryFixRequest, output *DictionaryFixResponse) error {
			sc, err := registry.Get(input.Code)
			if errors.Is(spellchecker.ErrNotFound, err) {
				return status.Wrap(err, status.NotFound)
			} else if err != nil {
				return status.Wrap(err, status.Internal)
			}

			fuzziness, err := input.Fuzziness.Parse()
			if err != nil {
				return status.Wrap(err, status.InvalidArgument)
			}

			if input.Text == "" {
				output.Fixes = make([]Fix, 0)
				return nil
			}

			matches := splitter.FindAllStringIndex(input.Text, -1)
			fixes := make([]Fix, 0, len(matches))
			correct := make([]Correct, 0, len(matches))

			for _, match := range matches {
				startByte, endByte := match[0], match[1]
				startRune := utf8.RuneCountInString(input.Text[:startByte])
				endRune := startRune + utf8.RuneCountInString(input.Text[startByte:endByte])

				fix := Fix{
					Start: startRune,
					End:   endRune,
				}

				word := input.Text[startByte:endByte]

				suggestions := sc.Suggest(&f1mspellchecker.SearchOptions{
					MaxErrors:  input.MaxErrors,
					FilterFunc: spellchecker.ScoringFunc(fuzziness, input.SimilarityThreshold),
				}, word, input.Limit)

				if suggestions.ExactMatch {
					correct = append(correct, Correct{
						Start: startRune,
						End:   endRune,
					})

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
			output.Correct = correct

			return nil
		},
	)

	u.SetTitle("Fix text")
	u.SetDescription(
		"Performs spellchecking on the given input text. Returns misspelled words along with suggested corrections, up to the specified limit per word.",
	)
	u.SetExpectedErrors(status.Internal, status.NotFound)

	return u
}

type FuzzinessValue string

// Parse converts FuzzinessValue into spellchecker.Fuzziness.
func (fv FuzzinessValue) Parse() (spellchecker.Fuzziness, error) {
	raw := strings.TrimSpace(strings.ToUpper(string(fv)))

	switch {
	case raw == "", raw == "AUTO":
		return spellchecker.AutoFuzziness{Low: 3, High: 6}, nil
	case strings.HasPrefix(raw, "AUTO:"):
		parts := strings.Split(strings.TrimPrefix(raw, "AUTO:"), ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid AUTO fuzziness format: %q", raw)
		}
		low, err1 := strconv.Atoi(parts[0])
		high, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid AUTO fuzziness values: %q", raw)
		}
		return spellchecker.AutoFuzziness{Low: low, High: high}, nil
	default:
		if n, err := strconv.Atoi(raw); err == nil {
			return spellchecker.FixedFuzziness(n), nil
		}
		return nil, fmt.Errorf("unknown fuzziness value: %q", raw)
	}
}

// UnmarshalJSON validates that fuzziness is passed as a JSON string.
func (fv *FuzzinessValue) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("fuzziness must be a string: %w", err)
	}
	*fv = FuzzinessValue(raw)
	return nil
}
