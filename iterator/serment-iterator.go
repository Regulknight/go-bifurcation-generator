package iterator

type SegmentIterator struct {
	start float64
	stop  float64
	step  float64
}

func NewSegmentIterator(start, stop, step float64) *SegmentIterator {
	return &SegmentIterator{
		start: start,
		stop:  stop,
		step:  step,
	}
}

func (iterator *SegmentIterator) GetIteratorChannel() <-chan float64 {
	out := make(chan float64)

	go func() {
		for i := iterator.start; i < iterator.stop; i += iterator.step {
			out <- i
		}

		close(out)
	}()

	return out
}
