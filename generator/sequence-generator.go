package generator

const maximumSequenceSize = 50000

func GetSequenceGenerator(valueGenerator <-chan float64) <-chan []float64 {
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
