package feature

type Feature struct {
	Type       string
	Location   Location
	Qualifiers []Qualifier
}

func NewFeature(t, l string) *Feature {
	var f Feature
	f.Type = t
	f.Location = *NewLocationFromString(l)
	f.Qualifiers = make([]Qualifier, 0)
	return &f
}

func (f *Feature) AddQualifier(t, v string) {
	q := NewQualifier(t, v)
	f.Qualifiers = append(f.Qualifiers, *q)
}
