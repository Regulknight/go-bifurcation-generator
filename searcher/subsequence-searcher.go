package searcher

import (
	"bifurcation-generator/searcher/comparator"
	"log"
)

const cycleAcceptCriteria int = 5

func isSubSequence(sequence, subsequence []float64) bool {
	//Тут у нас есть простор для оптимизации и мы можем выкинуть i=0 потому что наша подпоследовательность всегда берётся как слайс из конца последовательности
	//Но я посчитал это странным поведением для метода
	for i := 0; i < cycleAcceptCriteria; i++ {
		sequenceComparePart := sequence[len(sequence)-len(subsequence):]

		compResult := comparator.IsEqualsSlices(subsequence, sequenceComparePart)

		if !compResult {
			return false
		}

		sequence = sequence[:len(sequence)-len(subsequence)]
	}

	return true
}

// Возвращает nil если в последовательности нет подпоследовательностей
func IsContainsSubsequences(sequence []float64) []float64 {
	for i := 1; i < len(sequence)/cycleAcceptCriteria; i++ {
		subsequence := sequence[len(sequence)-i:]
		result := isSubSequence(sequence, subsequence)

		if result {
			return subsequence
		}
	}

	return nil
}

func GetCycle(sequenceChan <-chan []float64) []float64 {
	for sequence := range sequenceChan {
		searchResult := IsContainsSubsequences(sequence)
		if searchResult != nil {
			return searchResult
		}

		if len(sequence)%100 == 0 {
			log.Println(len(sequence))
		}
	}

	return nil
}
