package generator

const maximumSequenceSize = 50000

func getNextBifurcationValue(x, r float64) float64 {
	return r * x * (1 - x)
}

func getBifurcationGenerator(x0, r float64) <-chan float64 {
	out := make(chan float64)

	x := x0
	go func() {
		for {
			x = getNextBifurcationValue(x, r)
			out <- x
		}
	}()

	return out
}

func getSequenceGenerator(valueGenerator <-chan float64) <-chan []float64 {
	out := make(chan []float64)
	arrayState := [maximumSequenceSize]float64{}

	go func() {
		for i := 0; i < maximumSequenceSize; i++ {
			currentValue := <-valueGenerator
			arrayState[i] = currentValue
			out <- arrayState[:i]
		}

		close(out)
	}()

	return out
}

func GetBifurcationSequenceGenerator(x0, r float64) <-chan []float64 {
	return getSequenceGenerator(getBifurcationGenerator(x0, r))
}
