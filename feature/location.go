package feature

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

	// Start must be lesser than End
	if sl.Start > sl.End {
		panic(fmt.Sprintf("Sub-location (%s) with a start greater than end, while coordinate must be relative the direct strand.", s))
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
	l.End = e
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
func NewLocationFromString(s string) *Location {
	var l Location
	l.Start = 0
	l.End = 0
	l.Strand = 1
	l.SubCount = 0
	l.RevComp = false

	// Check if location is in complement strand
	if regexp.MustCompile(`^complement`).MatchString(s) {
		l.RevComp = true
		l.Strand = -1
		s = s[11:(len(s) - 1)]
	}

	// Check if the location has multiple sub-locations
	if regexp.MustCompile(`^join`).MatchString(s) {
		// Delete the join instruction
		s = s[5:(len(s) - 1)]

		// Split sub-location and treat-it
		sl := strings.Split(s, ",")
		l.SubLocations = make([]SubLocation, len(sl))
		for sli, slv := range sl {
			l.SubLocations[sli] = *NewSubLocationFromString(slv)
		}
		l.SubCount = len(sl)

		// Update Start, End and Strand attributes
		l.UpdateSES()
	} else {
		l.SubLocations = make([]SubLocation, 1)
		l.SubLocations[0] = *NewSubLocationFromString(s)
		l.SubCount = 1
		l.Start = l.SubLocations[0].Start
		l.End = l.SubLocations[0].End
		if l.SubLocations[0].RevComp {
			l.Strand = l.Strand * -1
		}
	}
	return &l
}

// Update Start, End and Strand values
func (l *Location) UpdateSES() {
	start := l.SubLocations[0].Start
	end := l.SubLocations[0].End

	if l.SubCount > 1 {
		for i := 1; i < l.SubCount; i++ {
			if start > l.SubLocations[i].Start {
				start = l.SubLocations[i].Start
			}
			if end < l.SubLocations[i].End {
				end = l.SubLocations[i].End
			}
		}
	}

	l.Start = start
	l.End = end

	if l.RevComp {
		if l.SubLocations[0].RevComp {
			// Weird case: double reverse complement!
			l.Strand = 1
		} else {
			l.Strand = -1
		}
	} else {
		if l.SubLocations[0].RevComp {
			l.Strand = -1
		} else {
			l.Strand = 1
		}
	}
}

func (l *Location) ToString() string {
	if l.SubCount > 0 {
		str := l.SubLocations[0].ToString()
		for i := 1; i < l.SubCount; i++ {
			str += "," + l.SubLocations[i].ToString()
		}

		if l.SubCount > 1 {
			str = "join(" + str + ")"
		}

		if l.RevComp {
			return "complement(" + str + ")"
		} else {
			return str
		}
	} else {
		// NOTE: Possible throw an error or a warning...
		return ""
	}
}

// Compute global length of the location
func (l *Location) Length() int {
	return l.End - l.Start + 1
}

// Compute spliced length of the location
func (l *Location) SplicedLength() int {
	if l.SubCount == 1 {
		return l.Length()
	} else {
		len := 0
		for _, sl := range l.SubLocations {
			len += sl.End - sl.Start + 1
		}
		return len
	}
}
