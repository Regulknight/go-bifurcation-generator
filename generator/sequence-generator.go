package generator

const MaximumSequenceSize = 50000

func GetSequenceGenerator(valueGenerator <-chan float64) <-chan []float64 {
	out := make(chan []float64)
	arrayState := [MaximumSequenceSize]float64{}

	go func() {
		for i := 0; i < MaximumSequenceSize; i++ {
			currentValue := <-valueGenerator
			arrayState[i] = currentValue
			out <- arrayState[:i]
		}

		close(out)
	}()

	return out
}
