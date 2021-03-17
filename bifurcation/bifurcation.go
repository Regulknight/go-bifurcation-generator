package bifurcation

import (
	"bifurcation-generator/generator"
)

type Bifurcation struct {
	BifurcationSequenceChannel <-chan []float64
	CurrentR                   float64
	CurrentX0                  float64
}

func GetBifurcationChannel(bifurcationParams *Params) <-chan *Bifurcation {
	out := make(chan *Bifurcation)

	go func() {
		rChan := bifurcationParams.rIterator.GetIteratorChannel()

		for r := range rChan {
			resultChannel := bifurcationParams.function.GetResultChannel(bifurcationParams.x0, r)
			sequenceGenerator := generator.GetSequenceGenerator(resultChannel)
			state := &Bifurcation{
				BifurcationSequenceChannel: sequenceGenerator,
				CurrentR:                   r,
				CurrentX0:                  bifurcationParams.x0,
			}
			out <- state
		}

		close(out)
	}()

	return out
}
