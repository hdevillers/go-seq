package feature

import "fmt"

type Qualifier struct {
	Tag   string
	Value Value
}

func NewQualifier(t, v string) (Qualifier, error) {
	q := new(Qualifier)
	var err error
	q.Tag = t
	q.Value, err = NewValue(v)

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
