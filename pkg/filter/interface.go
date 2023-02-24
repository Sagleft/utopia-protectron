package filter

type Filter interface {
	Use(message string) (isDetected bool)
}

func GetFilters() map[string]Filter {
	return map[string]Filter{
		"nil": NewInternalLinksFilter(),
		"nel": NewExternalLinksFilter(),
		"np":  NewNoPubkeyFilter(),
		"c":   NewChannelsFilter(),
	}
}
