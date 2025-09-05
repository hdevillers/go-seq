package feature

import (
	"fmt"
	"testing"
)

func TestCreateDefaultValue(t *testing.T) {
	v, err := NewValue()
	if err != nil {
		t.Errorf("Default constructor call should not return an error: %s", err)
	}
	if v.RawStr != D_RAWSTR {
		t.Error("Default 'RawStr' value is not set properly.")
	}
	if v.IsBool != D_ISBOOL {
		t.Error("Default 'IsBool' value is not set properly.")
	}
	if v.HasQuote != D_HASQUOTE {
		t.Error(("Default 'HasQuote' value is not set properly."))
	}
}

func TestCreateGivenValue(t *testing.T) {
	expect := "my_value"
	v, err := NewValue(expect)
	if err != nil {
		t.Errorf("This call to NewValue should not return an error: %s", err)
	}
	if v.RawStr != expect {
		t.Error("'RawStr' value is not set properly compare to a given value.")
	}
	if v.IsBool != D_ISBOOL {
		t.Error("Default 'IsBool' value is not set properly.")
	}
	if v.HasQuote != D_HASQUOTE {
		t.Error(("Default 'HasQuote' value is not set properly."))
	}
}

func TestTooManyValues(t *testing.T) {
	_, err := NewValue("value1", "value2")
	if err == nil {
		t.Error("Value create should have failed due to a too high number of inputs.")
	}
}

func TestToStringWithoutQuote(t *testing.T) {
	val := "Hello guys!"
	v, err := NewValue(val)
	if err != nil {
		t.Errorf("This call to NewValue should not return an error: %s", err)
	}

	// It is expected that v.HasQuote is false
	if v.HasQuote {
		t.Error("HasQuote attribute should be false.")
	}

	out := v.ToString()
	if out != val {
		t.Errorf("Stored value is different from the original one. Expected %s and obtained %s", val, out)
	}
}

func TestToStringWithQuote(t *testing.T) {
	val := "\"Hello guys!\""
	v, err := NewValue(val)
	if err != nil {
		t.Errorf("This call to NewValue should not return an error: %s", err)
	}

	// It is expected that v.HasQuote is true
	if !v.HasQuote {
		t.Error("HasQuote attribute should be true.")
	}

	out := v.ToString()
	if out != val {
		t.Errorf("Stored value is different from the original one. Expected %s and obtained %s", val, out)
	}
}

func TestToStringMutateQuote(t *testing.T) {
	val1 := "Hello guys!"
	val2 := fmt.Sprintf("\"%s\"", val1)
	v, err := NewValue(val1)
	if err != nil {
		t.Errorf("This call to NewValue should not return an error: %s", err)
	}

	// It is expected that v.HasQuote is false
	if v.HasQuote {
		t.Error("HasQuote attribute should be false.")
	}

	// Turn it into quoted value
	v.HasQuote = true

	out := v.ToString()
	if out != val2 {
		t.Errorf("Stored value is different from the original one. Expected %s and obtained %s", val2, out)
	}
}

func TestToStringMutateBool(t *testing.T) {
	val := "Hello guys!"
	v, err := NewValue(val)
	if err != nil {
		t.Errorf("This call to NewValue should not return an error: %s", err)
	}

	// It is expected that v.HasQuote is false
	if v.HasQuote {
		t.Error("HasQuote attribute should be false.")
	}

	// Turn it into boolean value
	v.IsBool = true

	out := v.ToString()
	if out != "" {
		t.Errorf("Expected an empty string, received %s.", out)
	}
}
