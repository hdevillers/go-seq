package seqio

import (
	"testing"
)

// Parse a regular EMBL (without reference in the header)
func TestEmblHeaderNoRef(t *testing.T) {
	seqIn := NewReader("../examples/EMBL/normal_noRef.embl", "embl", false)

	// Check if the reader was properly initiated
	err := seqIn.GetError()
	if err != nil {
		t.Errorf("Reader initiation failed: %s.", err)
	}
	defer seqIn.Close()

	// Parse the sequence
	if seqIn.Next() {
		// Check parsing error
		err = seqIn.GetError()
		if err != nil {
			t.Errorf("Sequence parsing failed: %s.", err)
		}

		// Define some comparison variables
		id := "X56734"

		// Check ID value
		seq := seqIn.seq
		if seq.Id != id {
			t.Errorf("Failed to parse sequence id, expected %s and retrieved %s.", id, seq.Id)
		}
	} else {
		// Should not returned false
		t.Error("Failed to call next sequence.")
	}
}
