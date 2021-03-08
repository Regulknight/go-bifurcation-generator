package bifurcation

type BifurcationFunction struct {
	Function       func(float64, float64) float64
	FunctionAsText string
}

func NewBifurcationFunction(bifurcationFunction func(float64, float64) float64, functionAsText string) *BifurcationFunction {
	return &BifurcationFunction{
		Function:       bifurcationFunction,
		FunctionAsText: functionAsText,
	}
}

func (generator *BifurcationFunction) GetResultChannel(x0, r float64) <-chan float64 {
	return getBifurcationFunctionChannel(generator.Function, x0, r)
}

func getBifurcationFunctionChannel(bifurcationFunction func(float64, float64) float64, x0, r float64) <-chan float64 {
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

func DefaultBifurcationFunction() *BifurcationFunction {
	return NewBifurcationFunction(defaultBifurcationFunction, "Xn+1 = r * Xn * (1 - Xn)")
}

func defaultBifurcationFunction(x, r float64) float64 {
	return r * x * (1 - x)
}
