package feature

import (
	"fmt"
	"regexp"
	"strconv"
)

type SubLocation struct {
	Start   int
	End     int
	SglBase bool
	UnkStr  bool
	UnkEnd  bool
	RevComp bool
	Between bool
	UnkLoc  bool
}

// Create a simple sublocation from coordinate pairs
func NewSubLocation(s, e int) *SubLocation {
	return &SubLocation{
		s, e, false, false,
		false, false, false,
		false,
	}
}

// Create a sublocation from a string
// (see insdc rules: https://www.insdc.org/documents/feature-table#3.4)
/* Supported format:
122			Single base location
122..330	Regular coordinate pairs
<122..330	Start boundary is unknown (larger)
122..>330	End boundary is unknown (larger)
122.330		Exact coordinate are unknown but included in the range
122^123     Point a position between the bases
*/
func NewSubLocationFromString(s string) *SubLocation {
	var sl SubLocation
	sl.SglBase = false
	sl.RevComp = false
	sl.UnkStr = false
	sl.UnkEnd = false
	sl.UnkLoc = false
	sl.Between = false

	// Check if sub location is in complement strand
	if regexp.MustCompile(`^complement`).MatchString(s) {
		sl.RevComp = true
		s = s[11:(len(s) - 1)]
	}

	// Check if sub location contain only valid characters
	if regexp.MustCompile(`^[<>\d\.\^]+$`).MatchString(s) {
		// Most classical case: pair of coordinates
		re := regexp.MustCompile(`^(<*)(\d+)\.\.(>*)(\d+)$`)
		test := re.FindStringSubmatch(s)
		if len(test) == 5 {
			sl.Start, _ = strconv.Atoi(test[2]) // No error, regex selection!
			sl.End, _ = strconv.Atoi(test[4])
			if test[1] == "<" {
				sl.UnkStr = true
			}
			if test[3] == ">" {
				sl.UnkEnd = true
			}
		} else {
			// Test other solutions
			// Single base location
			re = regexp.MustCompile(`^(\d+)$`)
			test = re.FindStringSubmatch(s)
			if len(test) == 2 {
				sl.Start, _ = strconv.Atoi(test[1])
				sl.End = sl.Start
				sl.SglBase = true
			} else {
				// Unknown location
				re = regexp.MustCompile(`^(\d+)(\.|\^)(\d+)$`)
				test = re.FindStringSubmatch(s)
				if len(test) == 4 {
					sl.Start, _ = strconv.Atoi(test[1]) // No error, regex selection!
					sl.End, _ = strconv.Atoi(test[3])
					if test[2] == "." {
						sl.UnkLoc = true
					} else {
						sl.Between = true
					}
				} else {
					// Unconsistant format
					panic(fmt.Sprintf("The sub-location %s has an unsupported format.", s))
				}
			}
		}
	} else {
		panic(fmt.Sprintf("The sub-location %s contains non valid characters.", s))
	}

	return &sl
}

func (sl *SubLocation) ToString() string {
	var str string

	if sl.SglBase {
		str = fmt.Sprintf("%d", sl.Start)
	} else {
		sep := ".." // Coordinate separator
		if sl.Between {
			sep = "^"
		} else if sl.UnkLoc {
			sep = "."
		}

		us := "" // Unknown start indicator
		if sl.UnkStr {
			us = "<"
		}

		ue := "" // Unknown end indicator
		if sl.UnkEnd {
			ue = ">"
		}

		str = fmt.Sprintf("%s%d%s%s%d", us, sl.Start, sep, ue, sl.End)
	}

	if sl.RevComp {
		return "complement(" + str + ")"
	} else {
		return str
	}
}

type Location struct {
	Start        int
	End          int
	RevComp      bool
	Strand       int
	SubLocations []SubLocation
	SubCount     int
}

// Create a simple location with a start, an end and a strand
func NewLocationSimple(s, e int, rc bool) *Location {
	if s > e {
		panic("In location, the start is supposed to be smaller than the end.")
	}

	var l Location
	l.Start = s
	l.End = s
	if rc {
		l.Strand = -1
		l.RevComp = true
	} else {
		l.Strand = 1
		l.RevComp = false
	}
	l.SubLocations = append(l.SubLocations, *NewSubLocation(s, e))
	l.SubCount = 1

	return &l
}

// Create a location from a string
//(see insdc rules: https://www.insdc.org/documents/feature-table#3.4)
func NewLocationString(s string) *Location {
	var l Location

	return &l
}
