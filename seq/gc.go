package seq

import (
	"errors"
)

type Gc struct {
	Counter []int64
}

func NewGc() Gc {
	gc := Gc{}
	gc.Counter = make([]int64, 255)
	bases := []byte{'C', 'c', 'G', 'g'}

	for i := 0; i < len(bases); i++ {
		gc.Counter[int(bases[i])] = 1
	}

	return gc
}

func (gc *Gc) Count(s *Seq) (int64, error) {
	var val int64 = 0
	l := s.Length()
	for i := 0; i < l; i++ {
		val = val + gc.Counter[int(s.Sequence[i])]
	}
	return val, nil
}

func (gc *Gc) CountIndex(s *Seq, f int, t int) (int64, error) {
	var val int64 = 0
	l := s.Length()
	if t > l {
		return val, errors.New("End index is higher than the sequence length.")
	}
	if f > t {
		return val, errors.New("Start index is higher than end index.")
	}
	if f < 1 {
		return val, errors.New("Start index must be greater or equal to 1.")
	}

	for i := f - 1; i < t; i++ {
		val = val + gc.Counter[int(s.Sequence[i])]
	}
	return val, nil
}
