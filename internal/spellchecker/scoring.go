package spellchecker

import (
	"math"

	"github.com/agext/levenshtein"
	"github.com/f1monkey/spellchecker/v2"
)

func ScoringFunc(maxErrors int, similarityThreshold float64) spellchecker.FilterFunc {
	return func(src, candidate []rune, count uint) (float64, bool) {
		distance, prefixLen, suffixLen := levenshtein.Calculate(src, candidate, 0, 1, 1, 1)
		if distance > maxErrors {
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
