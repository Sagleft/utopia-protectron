package filter

import "regexp"

var AllFiltersMap = getFiltersMap()
var marchURLRegExp = regexp.MustCompile(`\bhttps?://\S+\b`)
var marchIdyllRegExp = regexp.MustCompile(`\butopia://\S+\b`)

func isContainsURL(input string) bool {
	return marchURLRegExp.MatchString(input)
}

func isIdyllURL(input string) bool {
	return marchIdyllRegExp.MatchString(input)
}

func GetFiltersArray() []Filter {
	return []Filter{
		NewInternalLinksFilter(),
		NewExternalLinksFilter(),
		NewNoPubkeyFilter(),
		NewChannelsFilter(),
	}
}

func getFiltersMap() map[string]Filter {
	m := map[string]Filter{}
	filters := GetFiltersArray()

	for _, f := range filters {
		m[f.GetTag()] = f
	}
	return m
}
