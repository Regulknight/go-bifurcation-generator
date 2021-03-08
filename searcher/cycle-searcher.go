package searcher

import (
	"bifurcation-generator/generator"
)

type CycleSearcher struct {
	sequenceGeneratorChannel <-chan <-chan []float64
}

func NewCycleSearcher(sequenceGeneratorChannel <-chan <-chan []float64) *CycleSearcher {
	return &CycleSearcher{sequenceGeneratorChannel: sequenceGeneratorChannel}
}

func GetCyclesChannel(sequenceGenerator <-chan []float64) <-chan []float64 {
	out := make(chan []float64)

	go func() {
		for sequence, sequenceOk := <-sequenceGenerator; sequenceOk; {
			subsequence := IsContainsSubsequences(sequence)

			if subsequence != nil || len(sequence) > generator.MaximumSequenceSize {
				out <- subsequence
				sequenceOk = false
			} else {
				sequence, sequenceOk = <-sequenceGenerator
			}
		}

		close(out)
	}()

	return out
}
