package main

import (
	"bifurcation-generator/generator"
	"bifurcation-generator/iterator"
	"bifurcation-generator/subsequencesearcher"
)

func getSequenceGeneratorChannel(rIterator *iterator.SegmentIterator, sequenceGenerator *generator.BifurcationGenerator) <-chan <-chan []float64 {
	out := make(chan (<-chan []float64))

	go func() {
		rChan := rIterator.GetIteratorChannel()

		for r, rOk := <-rChan; rOk; r, rOk = <-rChan {
			sequenceGenerator := generator.GetSequenceGenerator(sequenceGenerator.GetResultChannel(0.4, r))
			out <- sequenceGenerator
		}
	}()

	return out
}

func getBifurcationCyclesChannel() <-chan []float64 {
	out := make(chan []float64)

	go func() {
		sequenceGeneratorChannel := getSequenceGeneratorChannel(iterator.NewSegmentIterator(0.0, 3.9, 0.1), generator.DefaultBifurcationGenerator())

		for sequenceGenerator, sgOk := <-sequenceGeneratorChannel; sgOk; sequenceGenerator, sgOk = <-sequenceGeneratorChannel {
			for sequence, sequenceOk := <-sequenceGenerator; sequenceOk; {
				subsequence := subsequencesearcher.IsContainsSubsequences(sequence)

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
