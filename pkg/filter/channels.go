package filter

import "regexp"

var matchChannelIDRegExp = regexp.MustCompile(`\b[0-9a-fA-F]{32}\b`)

type channelsFilter struct{ baseFilter }

func NewChannelsFilter() Filter {
	return channelsFilter{
		baseFilter: baseFilter{tag: "nc", name: "no-channels"},
	}
}

func (f channelsFilter) Use(message string) bool {
	return matchChannelIDRegExp.MatchString(message)
}
