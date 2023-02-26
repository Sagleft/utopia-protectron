package filter

type stickersFilter struct{ baseFilter }

func NewStickersFilter() Filter {
	return stickersFilter{
		baseFilter: baseFilter{tag: "ns", name: "no-stickers"},
	}
}

func (f stickersFilter) Use(message string) bool {
	return false // TODO
}
