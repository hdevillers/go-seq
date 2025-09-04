package feature

type Qualifier struct {
	Tag   string
	Value Value
}

func NewQualifier(t, v string) *Qualifier {
	var q Qualifier
	q.Tag = t
	q.Value = *NewValue(v)

	return &q
}
