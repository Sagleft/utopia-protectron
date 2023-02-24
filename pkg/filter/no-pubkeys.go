package filter

import "regexp"

var matchPubkeyRegExp = regexp.MustCompile(`\b[0-9a-fA-F]{64}\b`)

type pubkeyFilter struct{ baseFilter }

func NewNoPubkeyFilter() Filter {
	return pubkeyFilter{
		baseFilter: baseFilter{tag: "np", name: "no-pubkeys"},
	}
}

func (f pubkeyFilter) Use(message string) bool {
	return matchPubkeyRegExp.MatchString(message)
}
