package filter

import "regexp"

var matchPubkeyRegExp = regexp.MustCompile(`\b[0-9a-fA-F]{64}\b`)

type pubkeyFilter struct{}

func NewNoPubkeyFilter() Filter {
	return pubkeyFilter{}
}

func (f pubkeyFilter) Use(message string) bool {
	return matchPubkeyRegExp.MatchString(message)
}
