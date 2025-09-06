package feature

import (
	"testing"
)

func TestSimpleQualifier(t *testing.T) {
	tag := "note"
	val := "great description"
	q, err := NewQualifier(tag, val)
	if err != nil {
		t.Errorf("This call to NewQualifier should not return an error: %s", err)
	}

	// Check if tag value is well stored
	if q.Tag != tag {
		t.Errorf("Stored tag is wrong, expected %s but found %s.", tag, q.Tag)
	}

	// Check if value is well stored
	if q.Value.RawStr != val {
		t.Errorf("Stored value is wrong, expected %s but found %s.", val, q.Value.RawStr)
	}

	// The provided value has no quotes
	if q.Value.HasQuote {
		t.Error("This qualifier should have no quotes.")
	}

	// This qualifier is not boolean
	if q.Value.IsBool {
		t.Error("This qualifier should not be boolean.")
	}
}

func TestBooleanQualifier(t *testing.T) {
	tag := "pseudo"
	q, err := NewQualifier(tag, "") // Boolean from empty string
	if err != nil {
		t.Errorf("This call to NewQualifier should not return an error: %s", err)
	}

	// Check if tag value is well stored
	if q.Tag != tag {
		t.Errorf("Stored tag is wrong, expected %s but found %s.", tag, q.Tag)
	}

	// This qualifier is boolean
	if !q.Value.IsBool {
		t.Error("This qualifier should not be boolean.")
	}

	// Get the qualifier string
	qs := q.ToString("/", "")
	ex := "/pseudo"
	if qs != ex {
		t.Errorf("Failed to get the right qualifier string, expected %s, obtained %s", ex, qs)
	}
}

func TestNoValueQualifier(t *testing.T) {
	tag := "pseudo"
	q, err := NewQualifier(tag) // Boolean from empty string
	if err != nil {
		t.Errorf("This call to NewQualifier should not return an error: %s", err)
	}

	// Check if tag value is well stored
	if q.Tag != tag {
		t.Errorf("Stored tag is wrong, expected %s but found %s.", tag, q.Tag)
	}

	// This qualifier is boolean
	if !q.Value.IsBool {
		t.Error("This qualifier should not be boolean.")
	}

	// Get the qualifier string
	qs := q.ToString("/", "")
	ex := "/pseudo"
	if qs != ex {
		t.Errorf("Failed to get the right qualifier string, expected %s, obtained %s", ex, qs)
	}
}

func TestValueWithQuoteQualifier(t *testing.T) {
	tag := "product"
	val := "\"hypothetical protein\""
	q, err := NewQualifier(tag, val) // Boolean from empty string
	if err != nil {
		t.Errorf("This call to NewQualifier should not return an error: %s", err)
	}

	// Check if tag value is well stored
	if q.Tag != tag {
		t.Errorf("Stored tag is wrong, expected %s but found %s.", tag, q.Tag)
	}

	// The provided value has quotes
	if !q.Value.HasQuote {
		t.Error("This qualifier should have quotes.")
	}

	// This qualifier is not boolean
	if q.Value.IsBool {
		t.Error("This qualifier should not be boolean.")
	}

	// Get the qualifier string
	qs := q.ToString("/", "=")
	ex := "/product=\"hypothetical protein\""
	if qs != ex {
		t.Errorf("Failed to get the right qualifier string, expected %s, obtained %s", ex, qs)
	}
}
