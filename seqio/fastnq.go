package seqio

import (
	"github.com/hdevillers/go-seq/seq"
)

// Fastq sequence reader struct
type FastnqReader struct {
	scan     FileScanner
	currId   string
	eof      bool
	waitQual bool
}

// Generate a new reader
func NewFastnqReader(fs FileScanner) *FastnqReader {
	return &FastnqReader{
		scan:     fs,
		currId:   "",
		eof:      false,
		waitQual: false,
	}
}

// Return true if reachs the end-of-file
func (r *FastnqReader) IsEOF() bool {
	return r.eof
}

// Read a single fastq entry
func (r *FastnqReader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var newSeq seq.Seq
	iseof := true

	/*
		NOTE: this parser version is made to save time:
		1) Quality is not considered
		2) Fastq file is supposed well formed
	*/

	if r.scan.Scan() {
		// Check possible scanning error
		err := r.scan.Err()
		if err != nil {
			return newSeq, err
		}
		iseof = false

		// Get the ID line
		line := r.scan.Bytes()
		newSeq.SetId(string(line[1:]))

		// Get the sequence line
		r.scan.Scan()
		line = r.scan.Bytes()
		newSeq.AppendSequence(line)

		// Skip spacer line and quality line
		r.scan.Scan()
		r.scan.Scan()
	}

	if iseof {
		r.eof = true
	}

	return newSeq, nil
}
