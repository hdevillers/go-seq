package seqio

import (
	"errors"
	"fmt"

	"github.com/hdevillers/go-seq/seq"
)

// EMBL sequence reader struct
type EmblReader struct {
	scan     FileScanner
	currId   string
	currDesc string
	eof      bool
}

// EMBL sequence write struct
type EmblWriter struct {
	write FileWriter
	Count int
}

// Generate a new reader
func NewEmblReader(fs FileScanner) *EmblReader {
	return &EmblReader{
		scan:     fs,
		currId:   "",
		currDesc: "",
		eof:      false,
	}
}

// Generate a new writer
func NewEmblWriter(fw FileWriter) *EmblWriter {
	return &EmblWriter{
		write: fw,
		Count: 0,
	}
}

// Check if reached end of file
func (r *EmblReader) IsEOF() bool {
	return r.eof
}

// EMBL Read method
func (r *EmblReader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var newSeq seq.Seq

	// First scan the header lines
	hasID := false // Control that entry has an ID line
	hasFH := false // Control that entry has an FH line
	lastTag := ""
HEADER:
	for r.scan.Scan() {
		// Check possible scanning error
		err := r.scan.Err()
		if err != nil {
			return newSeq, err
		}

		// Get the scanned line
		line := r.scan.Bytes()

		// Skip possible blank line
		for len(line) < 2 {
			ok := r.scan.Scan()
			if !ok {
				// Reach end of file, finding no data
				return newSeq, errors.New("the provided file ends with blank line(s) and several mandatory header line(s) are missing")
			}
			// Check scanning error
			err := r.scan.Err()
			if err != nil {
				return newSeq, err
			}
			// Reset line data
			line = r.scan.Bytes()
		}

		// get the line tag
		tag := string(line[0:2])

		// Parse the line according to the tag
		switch tag {
		case "XX":
			// Just skip separator lines (most frequent ones)
			continue
		case "ID":
			hasID = true // Found the ID line
			// parseIdLine(line)
		case "FH":
			hasFH = true // Found the FH line
			lastTag = "FH"
			break HEADER
		case "FT":
			// Should not be possible, must have found FH before
			lastTag = "FT"
			break HEADER // hasFH will be false
		case "SQ":
			// Entry does not have features
			lastTag = "SQ"
			break HEADER
		case "  ":
			// Supposed to be a sequence line, not possible
			return newSeq, errors.New("embl entry is malformed: missing tags while supposed to be in the header")
		case "//":
			// Reach end of entry, not possible
			return newSeq, errors.New("reached end of entry without detecting features or sequence data")
		default:
			// unsupported tag return an error
			return newSeq, fmt.Errorf("unsupported line tag: %s", tag)
		}

		// Check if ID line was encountered
		if !hasID {
			return newSeq, errors.New("missing ID line the provided entry")
		}

		// Check if stopped at feature header
		if !hasFH {
			// Then only accept to be at the SQ line
			if lastTag != "SQ" {
				// NOTE: could not be fatal, let see later...
				return newSeq, fmt.Errorf("embl entry is malformed: missing feature header line FH")
			}
		}

	}

	// Return the populated sequence object
	return newSeq, nil
}
