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

func TestValueToString(t *testing.T) {
	pre := "FT   "
	tag := "name"
	val := "blabla"
	nch := 20

	v := NewValue(val)
	out := v.ToString(tag, pre, nch)
	exl := len(pre) + len(tag) + len(val) + 4
	if len(out) != exl {
		t.Error("Value to string convertion failed (IsBool: false, HasQuote: true).")
	}

	v.HasQuote = false
	out = v.ToString(tag, pre, nch)
	exl = len(pre) + len(tag) + len(val) + 2
	if len(out) != exl {
		t.Error("Value to string convertion failed (IsBool: false, HasQuote: false).")
	}

	v.IsBool = true
	out = v.ToString(tag, pre, nch)
	exl = len(pre) + len(tag) + 1
	if len(out) != exl {
		t.Error("Value to string convertion failed (IsBool: true, HasQuote: false).")
	}
}

func TestTooSmallLine(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("ToString call should have faile due to too short line length.")
		}
	}()

	pre := "FT   "
	tag := "name"
	val := "blabla"
	nch := 8

	v := NewValue(val)
	_ = v.ToString(tag, pre, nch)
}
