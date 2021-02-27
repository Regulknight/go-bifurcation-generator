package main

import "fmt"
import "math"

func getNextBifurcationValue(x, r float64) float64 {
    return r * x * (1 - x)
}

func bifurcationFunctionGenerator(x0, r  float64) <- chan float64 {
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

const cycleAcceptCriteria int = 5
const cycleMaximumSize int = 10000

func generateArrayStates(valueGenerator <- chan float64) <- chan []float64 {
    out := make (chan []float64)
    
    const bufferSize = cycleAcceptCriteria * cycleMaximumSize
    arrayState := [bufferSize]float64{}

    go func() {
        for i := 0; i < bufferSize; i++ {
            currentValue := <- valueGenerator
            arrayState[i] = currentValue
            out <- arrayState[:i]
        }

        close(out)
    }()

    return out
}

const roundingErrorLimit float64 = 0.01

func isItEquals(firstValue, secondValue float64) bool {
    return math.Pow(math.Abs((math.Pow(firstValue, 2.0) - math.Pow(secondValue, 2.0))), 1.0/2.0) < roundingErrorLimit
}

func floatComparator(firstValueChan, secondValueChan <- chan float64) <- chan bool {
    out := make(chan bool)

    go func() {
        firstValue, firstFlag := <- firstValueChan
        secondValue, secondFlag := <- secondValueChan

        for firstFlag && secondFlag {
            out <- isItEquals(firstValue, secondValue)

            firstValue, firstFlag = <- firstValueChan
            secondValue, secondFlag = <- secondValueChan
        }
        close(out)
    }()

    return out
}

func isEqualsSlices(firstSlice, secondSlice []float64) bool {
    if len(firstSlice) != len(secondSlice) {
        return false
    }

    for i := 0; i < len(firstSlice); i++ {
        if !isItEquals(firstSlice[i], secondSlice[i]) {
            return false
        }
    }

    return true
}

func sliceComparator(firstSliceChan, secondSliceChan <- chan []float64) <- chan bool {
    out := make(chan bool)

    firstSlice, firstFlag := <- firstSliceChan
    secondSlice, secondFlag := <- secondSliceChan

    go func() {
        for firstFlag && secondFlag {
            firstSlice, firstFlag = <- firstSliceChan
            secondSlice, secondFlag = <- secondSliceChan

            out <- isEqualsSlices(firstSlice, secondSlice)
        }
        close(out)
    }()

    return out
}

func getSliceStaticChannel(slice []float64, count int) <- chan []float64 {
    out := make(chan []float64)

    go func() {
        for i := 0; i < count; i++{
            out <- slice
        }
        close(out)
    }()

    return out
}

func getSlicerChannel(slice [] float64, sliceSize int) <- chan []float64 {
    out := make(chan []float64)

    go func() {
        for i := 1; i <= cycleAcceptCriteria; i++ {
            sliceStartIndex := len(slice) - i * sliceSize
            out <- slice[sliceStartIndex: sliceStartIndex + sliceSize]
        } 
    }()

    return out
}

func isItSubSequence(sequence, subsequence []float64) bool {
    //Тут у нас есть простор для оптимизации и мы можем выкинуть i=0 потому что наша подпоследовательность всегда берётся как слайс из конца последовательности
    //Но я посчитал это странным поведением для метода
    for i := 0; i < cycleAcceptCriteria; i++ {
        compChan := sliceComparator(getSliceStaticChannel(subsequence, cycleAcceptCriteria), getSlicerChannel(sequence, len(subsequence)))    

        result := <- compChan

        if !result {
            return false
        }
    }

    return true
}

// Возвращает nil если в последовательности нет подпоследовательностей
func isItContainsSubsequences(sequence []float64) []float64 {
    for i := 1; i < len(sequence) / cycleAcceptCriteria; i++ {
        subsequence := sequence[len(sequence) - i:]
        result := isItSubSequence(sequence, subsequence)

        if result {
            return subsequence
        }
    }

    return nil
}

func main() {
    for r := 0.0; r < 3.9; r+= 0.1 {
        out := generateArrayStates(bifurcationFunctionGenerator(0.4, r))
    
            calculationSlice, ok := <- out
    
            for ok {
                subsequence := isItContainsSubsequences(calculationSlice)
    
            if subsequence != nil || len(calculationSlice) > 40000 {
                fmt.Println(subsequence)
                ok = false
            } else {
                calculationSlice, ok = <-out
            }
        }
    }

	fmt.Println("Thats all")
}