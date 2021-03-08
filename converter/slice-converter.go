package converter

import "strconv"

func GetFloatSliceConverter(float64Chan <-chan []float64) <-chan []byte {
	out := make(chan []byte)

	go func() {
		dst := []byte{}
		floatSlice := <-float64Chan
		for i := 0; i < len(floatSlice); i++ {
			dst = append(dst[:], strconv.AppendFloat(dst, floatSlice[i], 'E', -1, 32)[:]...)
		}

		out <- dst
	}()

	return out
}
