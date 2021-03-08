package comparator

import "math"

const roundingErrorLimit float64 = 0.01

func isEqualsFloats(firstValue, secondValue float64) bool {
	return math.Pow(math.Abs(math.Pow(firstValue, 2.0)-math.Pow(secondValue, 2.0)), 1.0/2.0) < roundingErrorLimit
}

func isEqualsSlices(firstSlice, secondSlice []float64) bool {
	if len(firstSlice) != len(secondSlice) {
		return false
	}

	for i := 0; i < len(firstSlice); i++ {
		if !isEqualsFloats(firstSlice[i], secondSlice[i]) {
			return false
		}
	}

	return true
}

func GetSliceComparator(firstSliceChan, secondSliceChan <-chan []float64) <-chan bool {
	out := make(chan bool)

	firstSlice, firstFlag := <-firstSliceChan
	secondSlice, secondFlag := <-secondSliceChan

	go func() {
		for firstFlag && secondFlag {
			firstSlice, firstFlag = <-firstSliceChan
			secondSlice, secondFlag = <-secondSliceChan

			out <- isEqualsSlices(firstSlice, secondSlice)
		}
		close(out)
	}()

	return out
}
