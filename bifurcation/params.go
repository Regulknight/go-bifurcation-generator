package bifurcation

import "bifurcation-generator/iterator"

type Params struct {
	function  *RecurrentFunction
	rIterator *iterator.SegmentIterator
	x0        float64
}

func NewParams(bifurcationFunction *RecurrentFunction, rIterator *iterator.SegmentIterator, x0 float64) *Params {
	return &Params{
		function:  bifurcationFunction,
		rIterator: rIterator,
		x0:        x0,
	}
}
