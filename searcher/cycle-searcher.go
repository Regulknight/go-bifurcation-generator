package searcher

import (
	"github.com/regulknight/go-bifurcation-generator/generator"
)

type CycleSearcher struct {
	sequenceGeneratorChannel <-chan <-chan []float64
}

func NewCycleSearcher(sequenceGeneratorChannel <-chan <-chan []float64) *CycleSearcher {
	return &CycleSearcher{sequenceGeneratorChannel: sequenceGeneratorChannel}
}

func (cs *CycleSearcher) GetCyclesChannel() <-chan []float64 {
	out := make(chan []float64)

	go func() {
		for sequenceGenerator, sgOk := <-cs.sequenceGeneratorChannel; sgOk; sequenceGenerator, sgOk = <-cs.sequenceGeneratorChannel {
			for sequence, sequenceOk := <-sequenceGenerator; sequenceOk; {
				subsequence := IsContainsSubsequences(sequence)

				if subsequence != nil || len(sequence) > generator.MaximumSequenceSize {
					out <- subsequence
					sequenceOk = false
				} else {
					sequence, sequenceOk = <-sequenceGenerator
				}
			}
		}

	}()

	return out
}
