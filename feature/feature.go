package feature

type Feature struct {
	Type       string
	Location   Location
	Qualifiers map[string]Value
}

func NewFeature(t, l string) *Feature {
	var f Feature
	f.Type = t
	f.Location = *NewLocationFromString(l)
	return &f
}
