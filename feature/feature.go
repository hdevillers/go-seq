package feature

type Feature struct {
	Type       string
	Location   Location
	Qualifiers []Qualifier
}

func NewFeature(t, l string) *Feature {
	f := new(Feature)
	f.Type = t
	f.Location = *NewLocationFromString(l)
	f.Qualifiers = make([]Qualifier, 0)
	return f
}

func (f *Feature) AddQualifier(t, v string) error {
	q, err := NewQualifier(t, v)
	if err != nil {
		return err
	}
	f.Qualifiers = append(f.Qualifiers, q)
	return nil
}
