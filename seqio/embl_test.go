package seqio

import (
	"testing"
)

func compareValues(t *testing.T, v string, a *map[string][]string, w string) {
	tab, ok := (*a)[w]
	if !ok {
		// Key is missing
		t.Errorf("No '%s' entry in the annotation data.", w)
	} else {
		// Testing simple value, then expecting one entry
		if len(tab) != 1 {
			t.Errorf("Expected one entry in annotation for '%s', found %d.", w, len(tab))
		} else {
			if tab[0] != v {
				t.Errorf("Stored value in annotation data for '%s' is wrong. Expected %s, found %s.", w, v, tab[0])
			}
		}
	}
}

func compareArrays(t *testing.T, v *[]string, a *map[string][]string, w string) {
	tab, ok := (*a)[w]
	if !ok {
		// Key is missing
		t.Errorf("No '%s' entry in the annotation data.", w)
	} else {
		// array should have the same length
		if len(tab) != len(*v) {
			t.Errorf("Expected %d entries in annotation for '%s', found %d.", len(*v), w, len(tab))
		} else {
			for i, j := range *v {
				if tab[i] != j {
					t.Errorf("Value #%d is wrong in annotation data for '%s' is wrong. Expected %s, found %s.", i, w, j, tab[i])
				}
			}
		}
	}
}

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
		sv := "1"
		top := "linear"
		mty := "mRNA"
		cls := "STD"
		org := "PLN"
		len := "1859"
		ac := []string{"X56734", "S46826"}
		dt1 := "12-SEP-1991"
		dt2 := "25-NOV-2005"
		rl1 := "29"
		rl2 := "85"
		de := "Trifolium repens mRNA for non-cyanogenic beta-glucosidase"
		kw := []string{"beta-glucosidase"}
		os := "Trifolium repens (white clover)"
		oc := []string{"Eukaryota", "Viridiplantae", "Streptophyta", "Embryophyta", "Tracheophyta", "Spermatophyta", "Magnoliophyta", "eudicotyledons", "core eudicotyledons", "rosids", "fabids", "Fabales", "Fabaceae", "Papilionoideae", "Trifolieae", "Trifolium"}
		dr := []string{"EuropePMC", "PMC99098", "11752244"}

		// Check ID value
		seq := seqIn.seq
		if seq.Id != id {
			t.Errorf("Failed to parse sequence id, expected %s and retrieved %s.", id, seq.Id)
		}

		// Check description
		if seq.Desc != de {
			t.Errorf("Failed to parse description, expected %s and retrieved %s.", de, seq.Desc)
		}

		// Annotation content
		compareValues(t, sv, &seq.Annotations, "version")
		compareValues(t, top, &seq.Annotations, "topology")
		compareValues(t, mty, &seq.Annotations, "molecular_type")
		compareValues(t, cls, &seq.Annotations, "data_class")
		compareValues(t, org, &seq.Annotations, "taxonomic_division")
		compareValues(t, len, &seq.Annotations, "sequence_length")
		compareValues(t, dt1, &seq.Annotations, "created_date")
		compareValues(t, dt2, &seq.Annotations, "last_date")
		compareValues(t, rl1, &seq.Annotations, "created_release")
		compareValues(t, rl2, &seq.Annotations, "last_release")
		compareValues(t, os, &seq.Annotations, "species")
		compareArrays(t, &ac, &seq.Annotations, "accession")
		compareArrays(t, &kw, &seq.Annotations, "keywords")
		compareArrays(t, &oc, &seq.Annotations, "classification")
		compareArrays(t, &dr, &seq.Annotations, "cross_reference")

	} else {
		// Should not returned false
		t.Error("Failed to call next sequence.")
	}
}

// Parse a regular EMBL (without reference in the header)
func TestEmblFeatures(t *testing.T) {
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

		// Load the sequence
		seq := seqIn.seq

		// The testing file contain 3 features
		if len(seq.Features) != 3 {
			t.Errorf("The embl entry should contain 3 features, found %d.", len(seq.Features))
		}

		// The first one should be source
		if seq.Features[0].Type != "source" {
			t.Errorf("The first feature type should be 'source', found %s.", seq.Features[0].Type)
		}

		// Source start
		if seq.Features[0].Location.Start != 1 {
			t.Errorf("Feature 'source' start should be 1, found %d.", seq.Features[0].Location.Start)
		}

		// Source end
		if seq.Features[0].Location.End != 1859 {
			t.Errorf("Feature 'source' end should be 1859, found %d.", seq.Features[0].Location.End)
		}

		// Source should have 6 qualifiers
		if len(seq.Features[0].Qualifiers) != 6 {
			t.Errorf("Feature 'source' should have 6 qualifier entries, found %d.", len(seq.Features[0].Qualifiers))
		}

		// The second one should be mRNA
		if seq.Features[1].Type != "mRNA" {
			t.Errorf("The first feature type should be 'mRNA', found %s.", seq.Features[1].Type)
		}

		// mRNA should have 1 qualifier
		if len(seq.Features[1].Qualifiers) != 1 {
			t.Errorf("Feature 'source' should have 1 qualifier entry, found %d.", len(seq.Features[1].Qualifiers))
		}

		// The third one should be CDS
		if seq.Features[2].Type != "CDS" {
			t.Errorf("The first feature type should be 'CDS', found %s.", seq.Features[2].Type)
		}

		// CDS should have 11 qualifiers
		if len(seq.Features[2].Qualifiers) != 11 {
			t.Errorf("Feature 'source' should have 11 qualifier entries, found %d.", len(seq.Features[2].Qualifiers))
		}

		// Checking qualifier values
		// * value with quotes and on a single line
		if seq.Features[0].Qualifiers[4].Value.ToString() != "\"leaves\"" {
			t.Errorf("Bad qualifier value, expected \"leaves\", found %s.", seq.Features[0].Qualifiers[4].Value.ToString())
		}
		// * value without quote on a single line
		if seq.Features[2].Qualifiers[3].Value.ToString() != "1" {
			t.Errorf("Bad qualifier value, expected 1, found %s.", seq.Features[2].Qualifiers[3].Value.ToString())
		}
		// * value with quotes but on two lines
		if seq.Features[1].Qualifiers[0].Value.ToString() != "\"experimental evidence, no additional details recorded\"" {
			t.Errorf("Bad qualifier value, expected \"experimental evidence, no additional details recorded\", found %s.", seq.Features[1].Qualifiers[0].Value.ToString())
		}
		// * boolean value
		if !seq.Features[2].Qualifiers[4].Value.IsBool {
			t.Error("Bad qualifier value, expected to be boolean, but stored as regular value.")
		}
	}
}
