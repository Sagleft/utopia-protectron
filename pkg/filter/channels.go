package filter

type channelsFilter struct{}

func NewChannelsFilter() Filter {
	return channelsFilter{}
}

func (f channelsFilter) Use(message string) bool {
	return len(message) == channelIDLength &&
		isHexadecimal(message)
}
