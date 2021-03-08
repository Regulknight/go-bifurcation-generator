package generator

type BifurcationGenerator struct {
	Function       func(float64, float64) float64
	FunctionAsText string
}

func NewBifurcationGenerator(bifurcationFunction func(float64, float64) float64, functionAsText string) *BifurcationGenerator {
	return &BifurcationGenerator{
		Function:       bifurcationFunction,
		FunctionAsText: functionAsText,
	}
}

func (generator *BifurcationGenerator) GetResultChannel(x0, r float64) <-chan float64 {
	return getBifurcationGeneratorChannel(generator.Function, x0, r)
}

func getBifurcationGeneratorChannel(bifurcationFunction func(float64, float64) float64, x0, r float64) <-chan float64 {
	out := make(chan float64)

	x := x0
	go func() {
		for {
			x = bifurcationFunction(x, r)
			out <- x
		}
	}()

	return out
}

func DefaultBifurcationGenerator() *BifurcationGenerator {
	return NewBifurcationGenerator(defaultBifurcationFunction, "Xn+1 = r * Xn * (1 - Xn)")
}

func defaultBifurcationFunction(x, r float64) float64 {
	return r * x * (1 - x)
}
