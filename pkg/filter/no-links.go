package filter

const channelIDLength = 32

type linksFilter struct{}

func NewLinksFilter() Filter {
	return linksFilter{}
}

func (f linksFilter) Use(message string) bool {
	return len(message) == channelIDLength &&
		isHexadecimal(message)
}
