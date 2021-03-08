package searcher

import (
	"bifurcation-generator/searcher/comparator"
)

const cycleAcceptCriteria int = 5

func getStaticSliceGenerator(slice []float64, count int) <-chan []float64 {
	out := make(chan []float64)

	go func() {
		for i := 0; i < count; i++ {
			out <- slice
		}
		close(out)
	}()

	return out
}

func getSlicerGenerator(slice []float64, sliceSize int) <-chan []float64 {
	out := make(chan []float64)

	go func() {
		for i := 1; i <= cycleAcceptCriteria; i++ {
			sliceStartIndex := len(slice) - i*sliceSize
			out <- slice[sliceStartIndex : sliceStartIndex+sliceSize]
		}
	}()

	return out
}
func isSubSequence(sequence, subsequence []float64) bool {
	//Тут у нас есть простор для оптимизации и мы можем выкинуть i=0 потому что наша подпоследовательность всегда берётся как слайс из конца последовательности
	//Но я посчитал это странным поведением для метода
	for i := 0; i < cycleAcceptCriteria; i++ {
		compChan := comparator.GetSliceComparator(getStaticSliceGenerator(subsequence, cycleAcceptCriteria), getSlicerGenerator(sequence, len(subsequence)))

		result := <-compChan

		if !result {
			return false
		}
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
