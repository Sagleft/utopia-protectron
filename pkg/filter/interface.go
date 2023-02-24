package filter

type Filter interface {
	Use(message string) (isDetected bool)
	GetTag() string
	GetName() string
}

func GetFiltersArray() []Filter {
	return []Filter{
		NewInternalLinksFilter(),
		NewExternalLinksFilter(),
		NewNoPubkeyFilter(),
		NewChannelsFilter(),
	}
}

func GetFiltersMap() map[string]Filter {
	m := map[string]Filter{}
	filters := GetFiltersArray()

	for _, f := range filters {
		m[f.GetTag()] = f
	}
	return m
}
