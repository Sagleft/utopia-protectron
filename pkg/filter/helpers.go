package filter

import "regexp"

var matchHexRegExp = regexp.MustCompile(`^[0-9a-fA-F]+$`)

func isHexadecimal(input string) bool {
	return matchHexRegExp.MatchString(input)
}
