package feature

type SubLocation struct {
	Start   int
	End     int
	ExtStr  bool
	ExtEnd  bool
	RevComp bool
	Between bool
	Unknown bool
}

// Create a simple sublocation from coordinate pairs
func NewSubLocation(s, e int) *SubLocation {
	return &SubLocation{
		s, e, false, false,
		false, false, false,
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

	return &sl
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
