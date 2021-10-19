package utils

import (
	"regexp"

	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio"
)

/*
	Load sequences in a ref array
*/
func LoadSeqInArray(i, f string, a *[]seq.Seq) int {
	nseq := 0
	isgz := false
	if regexp.MustCompile(`\.gz$`).MatchString(i) {
		isgz = true
	}
	reader := seqio.NewReader(i, f, isgz)
	reader.CheckPanic()
	defer reader.Close()

	for reader.Next() {
		reader.CheckPanic()
		*a = append(*a, reader.Seq())
		nseq++
	}

	return nseq
}

/*
	Load sequences in a ref hash (indexed by sequence ID)
*/
func LoadSeqInMap(i, f string, m *map[string]seq.Seq) int {
	nseq := 0
	isgz := false
	if regexp.MustCompile(`\.gz$`).MatchString(i) {
		isgz = true
	}
	reader := seqio.NewReader(i, f, isgz)
	reader.CheckPanic()
	defer reader.Close()

	for reader.Next() {
		reader.CheckPanic()
		s := reader.Seq()
		(*m)[s.Id] = s
		nseq++
	}
	return nseq
}
