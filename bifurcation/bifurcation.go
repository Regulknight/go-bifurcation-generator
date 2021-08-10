package bifurcation

import (
	"github.com/regulknight/go-bifurcation-generator/generator"
	"github.com/regulknight/go-bifurcation-generator/iterator"
)

type Bifurcation struct {
	function  *BifurcationFunction
	rIterator *iterator.SegmentIterator
	x0        float64
}

func NewBifurcation(bifurcationFunction *BifurcationFunction, rIterator *iterator.SegmentIterator, x0 float64) *Bifurcation {
	return &Bifurcation{
		function:  bifurcationFunction,
		rIterator: rIterator,
		x0:        x0,
	}
}

func (bifurcation *Bifurcation) GetBifurcationChannel() <-chan <-chan []float64 {
	out := make(chan (<-chan []float64))

	go func() {
		rChan := bifurcation.rIterator.GetIteratorChannel()

		for r, rOk := <-rChan; rOk; r, rOk = <-rChan {
			sequenceGenerator := generator.GetSequenceGenerator(bifurcation.function.GetResultChannel(bifurcation.x0, r))
			out <- sequenceGenerator
		}
	}()

	return out
}
