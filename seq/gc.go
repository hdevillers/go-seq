package seq

import (
	"errors"
)

type Gc struct {
	Counter []int64
	Skipper []int64
}

func NewGc() Gc {
	gc := Gc{}
	gc.Counter = make([]int64, 255)
	gc.Skipper = make([]int64, 255)
	bases := []byte{'C', 'c', 'G', 'g'}

	for i := 0; i < len(bases); i++ {
		gc.Counter[int(bases[i])] = 1
	}

	// TODO: manage degenerated IUPAC bases?

	// Skip N in sequences
	gc.Skipper[int('N')] = 1
	gc.Skipper[int('n')] = 1

	return gc
}

func (gc *Gc) Count(s *Seq) (int64, int64, error) {
	var val int64 = 0
	var skp int64 = 0
	l := s.Length()
	for i := 0; i < l; i++ {
		val += gc.Counter[int(s.Sequence[i])]
		skp += gc.Skipper[int(s.Sequence[i])]
	}
	return val, skp, nil
}

func (gc *Gc) CountIndex(s *Seq, f int, t int) (int64, int64, error) {
	var val int64 = 0
	var skp int64 = 0
	l := s.Length()
	if t > l {
		return val, skp, errors.New("End index is higher than the sequence length.")
	}
	if f > t {
		return val, skp, errors.New("Start index is higher than end index.")
	}
	if f < 1 {
		return val, skp, errors.New("Start index must be greater or equal to 1.")
	}

	for i := f - 1; i < t; i++ {
		val += gc.Counter[int(s.Sequence[i])]
		skp += gc.Skipper[int(s.Sequence[i])]
	}
	return val, skp, nil
}
