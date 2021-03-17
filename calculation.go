package main

import (
	"bifurcation-generator/bifurcation"
	"bifurcation-generator/iterator"
	"bifurcation-generator/searcher"
	"encoding/json"
	"fmt"
	"time"
)

type CalculationResult struct {
	Cycle      []float64
	X          float64
	R          float64
	I          int
	SearchTime int
}

func NewCalculationResult(cycle []float64, x, r float64, i int, searchTime int) *CalculationResult {
	return &CalculationResult{
		Cycle:      cycle,
		X:          x,
		R:          r,
		I:          i,
		SearchTime: searchTime,
	}
}

func getCalculationResultChannel() <-chan *CalculationResult {
	calculationResultChannel := make(chan *CalculationResult)

	go func() {
		bifurcationGenerator := bifurcation.GetBifurcationChannel(bifurcation.NewParams(bifurcation.DefaultRecurrentFunction(), iterator.NewSegmentIterator(0.0, 3.9, 0.001), 0.4))

		i := 0
		for currentBifurcation := range bifurcationGenerator {
			startTime := time.Now().Nanosecond()
			cycle := searcher.GetCycle(currentBifurcation.BifurcationSequenceChannel)
			searchTime := time.Now().Nanosecond() - startTime

			calculationResult := NewCalculationResult(cycle, currentBifurcation.CurrentX0, currentBifurcation.CurrentR, i, searchTime)
			calculationResultChannel <- calculationResult
			i += 1
		}

		close(calculationResultChannel)
	}()

	return calculationResultChannel
}

func (result *CalculationResult) convertToJSON() string {
	b, err := json.Marshal(result)

	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	return string(b)
}

func ResultChannelConvertToByteArrayChannel(resultChannel <-chan *CalculationResult) <-chan []byte {
	out := make(chan []byte)

	go func() {
		for result := range resultChannel {
			jsonResult := result.convertToJSON()
			out <- []byte(jsonResult)
		}

		close(out)
	}()

	return out
}
