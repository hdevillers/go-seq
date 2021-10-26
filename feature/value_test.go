package feature

import (
	"testing"
)

func TestCreateDefaultValue(t *testing.T) {
	v := NewValue()
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
	v := NewValue(expect)
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
	defer func() {
		if r := recover(); r == nil {
			t.Error("Value create should have failed due to a too high number of inputs.")
		}
	}()

	v := NewValue("value1", "value2")

	if !v.IsBool {
		t.Error("Useless error message, the script is supposed to panic before.")
	}
}
