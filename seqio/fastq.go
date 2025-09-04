package seqio

import (
	"errors"
	"fmt"
	"os"

	"github.com/hdevillers/go-seq/seq"
)

const (
	FastqIdPreffix byte = '@'
	FastqSpPreffix byte = '+'
)

// Fastq sequence reader struct
type FastqReader struct {
	scan     FileScanner
	currId   string
	eof      bool
	waitQual bool
}

// Fastq sequence writer struct
type FastqWriter struct {
	write FileWriter
	Count int
}

// Generate a new reader
func NewFastqReader(fs FileScanner) *FastqReader {
	return &FastqReader{
		scan:     fs,
		currId:   "",
		eof:      false,
		waitQual: false,
	}
}

// Generate a new writer
func NewFastqWriter(fw FileWriter) *FastqWriter {
	return &FastqWriter{
		write: fw,
		Count: 0,
	}
}

// Return true if reachs the end-of-file
func (r *FastqReader) IsEOF() bool {
	return r.eof
}

// Read a single fastq entry
func (r *FastqReader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var newSeq seq.Seq
	iseof := true

READ:
	for r.scan.Scan() {
		// Check possible scanning error
		err := r.scan.Err()
		if err != nil {
			return newSeq, err
		}
		iseof = false

		// Get the scanned line
		line := r.scan.Bytes()

		if len(line) > 0 {
			if line[0] == FastqIdPreffix {
				// This is an ID line
				newSeq.SetId(string(line[1:]))

				// Get the sequence line
				r.scan.Scan()
				if r.scan.Err() != nil {
					return newSeq, r.scan.Err()
				}
				line = r.scan.Bytes()
				newSeq.AppendSequence(line)

				// Skip the spacer line
				r.scan.Scan()
				if r.scan.Err() != nil {
					return newSeq, r.scan.Err()
				}

				// Get the quality line
				r.scan.Scan()
				if r.scan.Err() != nil {
					return newSeq, r.scan.Err()
				}
				line = r.scan.Bytes()
				qerr := newSeq.Quality.AppendStrScore(line)
				// qerr are not fatal, just thow it
				if qerr != nil {
					fmt.Fprintln(os.Stderr, qerr)
				}

				break READ
			}
		}
	}

	if iseof {
		r.eof = true
	}

	return newSeq, nil
}

func (w *FastqWriter) Write(s seq.Seq) error {
	// Check sequence validity
	if s.Id == "" {
		return errors.New("[FASTQ WRITER]: Missing sequence ID")
	}
	if s.Length() == 0 {
		return errors.New("[FASTQ WRITER]: Cannot write out empty sequences")
	}
	if len(s.Quality.StrScore) == 0 {
		// If the quality is empty, then generate a fake score
		s.Quality.GenerateDefaultScore(s.Length())
	}
	if s.Length() != len(s.Quality.StrScore) {
		return errors.New("[FASTQ WRITER]: Sequence and quality with different lengths")
	}

	// Add the ID
	_, err := w.write.Write([]byte{FastqIdPreffix})
	if err != nil {
		return err
	}
	w.write.Write([]byte(s.Id))
	w.write.Write([]byte{'\n'})

	// Add the sequence
	w.write.Write(s.Sequence)
	w.write.Write([]byte{'\n', FastqSpPreffix, '\n'})

	// Add the quaity
	w.write.Write(s.Quality.StrScore)
	w.write.Write([]byte{'\n'})

	w.Count++

	return err
}

func (w *FastqWriter) Flush() error {
	err := w.write.Flush()
	return err
}
