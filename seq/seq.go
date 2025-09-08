package seq

import (
	"github.com/hdevillers/go-seq/feature"
	"github.com/hdevillers/go-seq/quality"
)

type Seq struct {
	Id          string
	Desc        string
	Sequence    []byte
	Quality     quality.Quality
	Annotations map[string][]string
	Features    []*feature.Feature
}

func NewSeq(id string) *Seq {
	p := Seq{Id: id}
	p.Annotations = make(map[string][]string)
	p.Features = make([]*feature.Feature, 0)
	return &p
}

func (s *Seq) SetId(id string) {
	s.Id = id
}

func (s *Seq) SetDesc(desc string) {
	s.Desc = desc
}

func (s *Seq) SetSequence(sequence []byte) {
	s.Sequence = sequence
}

func (s *Seq) AppendSequence(sequence []byte) {
	s.Sequence = append(s.Sequence, sequence...)
}

func (s *Seq) Length() int {
	return len(s.Sequence)
}
