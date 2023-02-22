package filter

import "regexp"

var marchURLRegExp = regexp.MustCompile(`\bhttps?://\S+\b`)
var marchIdyllRegExp = regexp.MustCompile(`\butopia://\S+\b`)

func isContainsURL(input string) bool {
	return marchURLRegExp.MatchString(input)
}

func isIdyllURL(input string) bool {
	return marchIdyllRegExp.MatchString(input)
}
