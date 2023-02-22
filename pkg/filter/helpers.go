package filter

import "regexp"

var matchHexRegExp = regexp.MustCompile(`^[0-9a-fA-F]+$`)
var marchURLRegExp = regexp.MustCompile(`\bhttps?://\S+\b`)

func isHexadecimal(input string) bool {
	return matchHexRegExp.MatchString(input)
}

func isContainsURL(input string) bool {
	return marchURLRegExp.MatchString(input)
}
