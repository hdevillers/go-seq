package feature

type Feature struct {
	Type       string
	Location   Location
	Qualifiers []Qualifier
}

func NewFeature(t, l string) (*Feature, error) {
	f := new(Feature)
	var err error
	f.Type = t
	f.Location, err = NewLocationFromString(l)
	f.Qualifiers = make([]Qualifier, 0)
	return f, err
}

func (f *Feature) AddQualifier(t, v string) error {
	q, err := NewQualifier(t, v)
	if err != nil {
		return err
	}
	f.Qualifiers = append(f.Qualifiers, q)
	return nil
}
