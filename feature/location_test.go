package feature

import (
	"testing"
)

func TestCreateDefaultSubLocation(t *testing.T) {
	start := 500
	end := 1000
	sl := NewSubLocation(start, end)

	if sl.Start != start {
		t.Error("Start value is wrong in SubLocation.")
	}

	if sl.End != end {
		t.Error("End value is wrong in SubLocation.")
	}
}

func TestCreateSubLocationFromString(t *testing.T) {
	s := "500..1000"
	sl := NewSubLocationFromString(s)

	if sl.Start != 500 {
		t.Error("Failed to parse start value in SubLocation.")
	}

	if sl.End != 1000 {
		t.Error("Failed to parse end value in SubLocation.")
	}
}

func TestSubLocationReadWriteString(t *testing.T) {
	ss := []string{
		"500..1000", "complement(23..45)", "<10..34",
		"56..>1025", "1.2", "4875^5000", "complement(<1..>38)",
		"1", "58467542316548",
	}

	for _, s := range ss {
		sl := NewSubLocationFromString(s)
		sout := sl.ToString()
		if s != sout {
			t.Errorf("Failed to reproduice sub-location string. Expected: %s; Obtained: %s", s, sout)
		}
	}
}

func TestSubLoctionBadCharacter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Failed to detect a bad character in a string.")
		}
	}()

	s := "234..2O0"
	_ = NewSubLocationFromString(s)
}

func TestCreateLocationSimple(t *testing.T) {
	l := NewLocationSimple(100, 200, false)
	if l.Start != 100 {
		t.Errorf("Failed to retrieve location start. Expected: 100; Obtained: %d", l.Start)
	}
	if l.End != 200 {
		t.Errorf("Failed to retrieve location end. Expected: 200; Obtained: %d", l.End)
	}
	if l.RevComp {
		t.Error("The location was expected in direct strand, it is in reverse complement.")
	}
	if l.Strand != 1 {
		t.Errorf("Location strand is wrong. Expected: 1; Obtained: %d", l.Strand)
	}
	if l.SubCount != 1 {
		t.Errorf("Wrong number of sub-locations. Expected: 1; Obtained: %d", l.SubCount)
	}
}

func TestCreateLocationFromString(t *testing.T) {
	str := "complement(join(13..234,400..1000))"
	l := NewLocationFromString(str)

	if l.Start != 13 {
		t.Errorf("Failed to retrieve location start. Expected: 13; Obtained: %d", l.Start)
	}
	if l.End != 1000 {
		t.Errorf("Failed to retrieve location end. Expected: 1000; Obtained: %d", l.End)
	}
	if !l.RevComp {
		t.Error("The location was expected in reverse strand, it is in direct strand.")
	}
	if l.Strand == 1 {
		t.Errorf("Location strand is wrong. Expected: -1; Obtained: %d", l.Strand)
	}
	if l.SubCount != 2 {
		t.Errorf("Wrong number of sub-locations. Expected: 2; Obtained: %d", l.SubCount)
	}
}

func TestLocationReadWriteString(t *testing.T) {
	str := []string{
		"100..200",
		"complement(100..200)",
		"join(2..100,300..433)",
		"complement(join(13..234,400..1000))",
		"join(complement(400..1000),complement(13..234))",
		"complement(join(complement(400..1000),complement(13..234)))",
		"join(<1..200,300..400,500)",
		"<234..>678",
	}

	for _, s := range str {
		l := NewLocationFromString(s)
		ns := l.ToString()
		if s != ns {
			t.Errorf("Failed to read/write location string. Excpected: %s; Obtained: %s.", s, ns)
		}
	}
}
