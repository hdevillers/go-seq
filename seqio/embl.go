package seqio

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hdevillers/go-seq/feature"
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

/*
	HEADER READING METHODS
*/

// Parse ID line
func parseEmblIdLine(dt string, a *map[string][]string) (string, error) {
	// Initiate annotation entries linked to ID line
	(*a)["ID_line"] = make([]string, 0)
	(*a)["ID_msg"] = make([]string, 0)
	(*a)["version"] = make([]string, 0)
	(*a)["topology"] = make([]string, 0)
	(*a)["molecular_type"] = make([]string, 0)
	(*a)["data_class"] = make([]string, 0)
	(*a)["taxonomic_division"] = make([]string, 0)
	(*a)["sequence_length"] = make([]string, 0)

	// Split the provided string
	dts := strings.Split(dt, "; ")

	// Prepare regex
	reId := regexp.MustCompile(`^([\w_\-\.]+)`)

	// dts should be an array of 7 strings
	if len(dts) == 7 {
		// Check if ID is ok
		id := reId.FindStringSubmatch(dts[0])
		if len(id) == 0 {
			return "", fmt.Errorf("ID line is malformed: failed to find the sequence id: %s", dt)
		}

		// Get the sequence version
		tmp := regexp.MustCompile(`(\d+)$`).FindStringSubmatch(dts[1])
		if len(tmp) == 0 {
			(*a)["ID_line"] = append((*a)["ID_line"], dt)
			(*a)["ID_msg"] = append((*a)["ID_msg"], "ID line is malformed: failed to retrieve the sequence version")
		} else {
			(*a)["version"] = append((*a)["version"], tmp[0])
		}

		// Get the topology (no check done)
		(*a)["topology"] = append((*a)["topology"], dts[2])

		// Get the molecular type (no check done)
		(*a)["molecular_type"] = append((*a)["molecular_type"], dts[3])

		// Get the data class (no check done)
		(*a)["data_class"] = append((*a)["data_class"], dts[4])

		// Get the taxonomic division (no check done)
		(*a)["taxonomic_division"] = append((*a)["taxonomic_division"], dts[5])

		// Get the sequence length (just check it is numerical)
		tmp = regexp.MustCompile(`^(\d+)`).FindStringSubmatch(dts[6])
		if len(tmp) == 0 {
			(*a)["ID_line"] = append((*a)["ID_line"], dt)
			(*a)["ID_msg"] = append((*a)["ID_msg"], "ID line is malformed: failed to retrieve the sequence length")
		} else {
			(*a)["sequence_length"] = append((*a)["sequence_length"], tmp[0])
		}

		// Return the retrieved ID and no error
		return id[0], nil
	} else {
		// ID line is malformed
		// Try to get at least the ID of the sequence
		tmp := reId.FindStringSubmatch(dt)
		if len(tmp) == 1 {
			// Just keep the ID line as it is
			(*a)["ID_line"] = append((*a)["ID_line"], dt)
			(*a)["ID_msg"] = append((*a)["ID_msg"], "ID line is malformed: failed to find the seven element it should contain")
			return tmp[0], nil
		} else {
			// Failed to get at least an ID
			return "", fmt.Errorf("ID line is malformed: failed to find the sequence id: %s", dt)
		}
	}
}

// Parse value from line separated by semi-colon
func parseEmblSCLine(dt, key, end string, a *map[string][]string) {
	// Initiate annotation if required
	_, ok := (*a)[key]
	if !ok {
		(*a)[key] = make([]string, 0)
	}

	// Delete tailing semi-colon (if it exists)
	dt, _ = strings.CutSuffix(dt, end)

	// Replace space characters by nothing (just in case line misses some spaces)
	dt = strings.ReplaceAll(dt, "; ", ";")

	// Split data
	dts := strings.Split(dt, ";")

	// Store AC value(s)
	(*a)[key] = append((*a)[key], dts...)
}

// Parse date (DT) lines
func parseEmblDtLine(dt, a *map[string][]string) {
	// Initiate annotations (if required)
	_, ok := (*a)["DT_line"]
	if !ok {
		(*a)["DT_line"] = make([]string, 0)
		(*a)["DT_msg"] = make([]string, 0)
	}
}

// EMBL Read method
func (r *EmblReader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var newSeq seq.Seq

	// Initialize attributes
	newSeq.Annotations = make(map[string][]string)
	newSeq.Features = make([]*feature.Feature, 0)

	// First scan the header lines
	hasID := false // Control that entry has an ID line
	hasFH := false // Control that entry has an FH line
	lastTag := ""  // Keep last line tag found
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
			newSeq.Id, err = parseEmblIdLine(string(line[4:]), &newSeq.Annotations)
			if err != nil {
				return newSeq, err
			}
		case "AC":
			parseEmblSCLine(string(line[4:]), "accession", ";", &newSeq.Annotations)
		case "PR":
			parseEmblSCLine(string(line[4:]), "project_id", ";", &newSeq.Annotations)
		case "DT":

		case "DE":

		case "KW":
			parseEmblSCLine(string(line[4:]), "keywords", ".", &newSeq.Annotations)
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
