package bifurcation

type RecurrentFunction struct {
	Function       func(float64, float64) float64
	FunctionAsText string
}

func NewRecurrentFunction(recurrentFunction func(float64, float64) float64, functionAsText string) *RecurrentFunction {
	return &RecurrentFunction{
		Function:       recurrentFunction,
		FunctionAsText: functionAsText,
	}
}

func (generator *RecurrentFunction) GetResultChannel(x0, r float64) <-chan float64 {
	return getRecurrentFunctionChannel(generator.Function, x0, r)
}

func getRecurrentFunctionChannel(RecurrentFunction func(float64, float64) float64, x0, r float64) <-chan float64 {
	out := make(chan float64)

	x := x0
	go func() {
		for {
			x = RecurrentFunction(x, r)
			out <- x
		}
	}()

	return out
}

func DefaultRecurrentFunction() *RecurrentFunction {
	return NewRecurrentFunction(defaultRecurrentFunction, "Xn+1 = r * Xn * (1 - Xn)")
}

func defaultRecurrentFunction(x, r float64) float64 {
	return r * x * (1 - x)
}
