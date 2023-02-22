package filter

import "regexp"

var matchChannelIDRegExp = regexp.MustCompile(`\b[0-9a-fA-F]{32}\b`)

type channelsFilter struct{}

func NewChannelsFilter() Filter {
	return channelsFilter{}
}

func (f channelsFilter) Use(message string) bool {
	return matchChannelIDRegExp.MatchString(message)
}
