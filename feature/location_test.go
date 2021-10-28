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
