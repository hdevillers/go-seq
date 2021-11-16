package fastq

import (
	"errors"
	"fmt"
	"os"

	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio/seqitf"
)

const (
	IdPreffix byte = '@'
	SpPreffix byte = '+'
)

// Fastq sequence reader struct
type Reader struct {
	scan     seqitf.FileScanner
	currId   string
	eof      bool
	waitQual bool
}

// Fastq sequence writer struct
type Writer struct {
	write seqitf.FileWriter
	Count int
}

// Generate a new reader
func NewReader(fs seqitf.FileScanner) *Reader {
	return &Reader{
		scan:     fs,
		currId:   "",
		eof:      false,
		waitQual: false,
	}
}

// Generate a new writer
func NewWriter(fw seqitf.FileWriter) *Writer {
	return &Writer{
		write: fw,
		Count: 0,
	}
}

// Return true if reachs the end-of-file
func (r *Reader) IsEOF() bool {
	return r.eof
}

// Read a single fastq entry
func (r *Reader) Read() (seq.Seq, error) {
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
			if line[0] == IdPreffix {
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

/*
			if r.currId != "" {
				// Return the current sequence if not nil
				if newSeq.Length() == 0 {
					// Empty sequence or bad format
					return newSeq, errors.New("[FASTQ READER]: Empty sequence or bad format.")
				}

				// Seq sequence data
				newSeq.SetId(r.currId)

				// Save the new ID
				r.currId = string(line[1:]) // Only skip the line preffix
				r.waitQual = false

				// Return the sequence
				return newSeq, nil
			} else {
				// Save the new ID
				r.currId = string(line[1:])

				// The sequence object should be empty
				if newSeq.Length() > 0 {
					return newSeq, errors.New("[FASTQ READER]: Sequence without ID ou bad format.")
				}

				// Continue
			}
		} else {
			if line[0] == SpPreffix {
				// Finished to read sequence line(s)
				// Start reading the quality
				r.waitQual = true

				// At that step, newSeq.Length must not be null
				if newSeq.Length() == 0 {
					return newSeq, errors.New("[FASTQ READER]: Empty sequence or bad format.")
				}

				// Continue
			} else {
				// Read sequence data or quality data
				// NOTE: We accept non standard fastq with sequence on multiple lines
				if r.waitQual {
					qerr := newSeq.Quality.AppendStrScore(line)
					// qerr are not fatal, just thow it
					fmt.Fprintln(os.Stderr, qerr)
				} else {
					newSeq.AppendSequence(line)
				}
			}
		}
	}
	// Scanning is finished
	r.eof = true

	// Set last sequence ID
	newSeq.SetId(r.currId)

	if newSeq.Length() == 0 {
		return newSeq, errors.New("[FASTQ READER]: Last sequence is null.")
	}

	// Return with no error
	return newSeq, nil
}
*/

func (w *Writer) Write(s seq.Seq) error {
	// Check sequence validity
	if s.Id == "" {
		return errors.New("[FASTQ WRITER]: Missing sequence ID.")
	}
	if s.Length() == 0 {
		return errors.New("[FASTQ WRITER]: Cannot write out empty sequences.")
	}
	if len(s.Quality.StrScore) == 0 {
		// If the quality is empty, then generate a fake score
		s.Quality.GenerateDefaultScore(s.Length())
	}
	if s.Length() != len(s.Quality.StrScore) {
		return errors.New("[FASTQ WRITER]: Sequence and quality with different lengths.")
	}

	// Add the ID
	_, err := w.write.Write([]byte{IdPreffix})
	_, err = w.write.Write([]byte(s.Id))
	if err != nil {
		return err
	}
	_, err = w.write.Write([]byte{'\n'})

	// Add the sequence
	_, err = w.write.Write(s.Sequence)
	if err != nil {
		return err
	}
	_, err = w.write.Write([]byte{'\n', SpPreffix, '\n'})

	// Add the quaity
	_, err = w.write.Write(s.Quality.StrScore)
	if err != nil {
		return err
	}
	_, err = w.write.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	//err = w.write.Flush()
	w.Count++

	return err
}

func (w *Writer) Flush() error {
	err := w.write.Flush()
	return err
}
