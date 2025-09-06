package feature

import (
	"errors"
	"fmt"
)

type Qualifier struct {
	Tag   string
	Value Value
}

func NewQualifier(t string, v ...string) (Qualifier, error) {
	q := new(Qualifier)
	var err error
	q.Tag = t
	if len(v) == 0 {
		q.Value, err = NewValue()
	} else if len(v) == 1 {
		q.Value, err = NewValue(v[0])
	} else {
		err = errors.New("only one string object per qualifier value is allowed")
	}

	return *q, err
}

func (q *Qualifier) ToString(prefix, suffix string) string {
	// If the qualifier is boolean or the value is D_RAWSTR
	if q.Value.IsBool || q.Value.RawStr == D_RAWSTR {
		return fmt.Sprintf("%s%s", prefix, q.Tag)
	}

	// Else get the complete value
	return fmt.Sprintf("%s%s%s%s", prefix, q.Tag, suffix, q.Value.ToString())
}
