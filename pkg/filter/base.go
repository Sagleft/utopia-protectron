package filter

type baseFilter struct {
	name string
	tag  string
}

func (f baseFilter) GetTag() string {
	return f.tag
}

func (f baseFilter) GetName() string {
	return f.name
}
