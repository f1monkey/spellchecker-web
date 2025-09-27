package spellchecker

import (
	"math"

	"github.com/agext/levenshtein"
	"github.com/f1monkey/spellchecker/v2"
)

type Fuzziness interface {
	MaxAllowedErrors(wordLen int) int
}

type FixedFuzziness int

func (f FixedFuzziness) MaxAllowedErrors(_ int) int {
	return int(f)
}

type AutoFuzziness struct {
	Low, High int
}

func (a AutoFuzziness) MaxAllowedErrors(wordLen int) int {
	if wordLen < a.Low {
		return 0
	}

	if wordLen < a.High {
		return 1
	}

	return 2
}

func ScoringFunc(fuzziness Fuzziness, similarityThreshold float64) spellchecker.FilterFunc {
	return func(src, candidate []rune, count uint) (float64, bool) {
		distance, prefixLen, suffixLen := levenshtein.Calculate(src, candidate, 0, 1, 1, 1)
		if distance > fuzziness.MaxAllowedErrors(len(src)) {
			return 0, false
		}

		if similarityThreshold > 0 {
			similarity := levenshtein.Match(string(src), string(candidate), nil)
			if similarity < similarityThreshold {
				return 0, false
			}
		}

		mult := math.Log1p(float64(count)) * math.Pow(1.5, float64(prefixLen+suffixLen))

		return 1 / (1 + float64(distance*distance)) * mult, true
	}
}
