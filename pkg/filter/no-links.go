package filter

type internalLinksFilter struct{}

type externalLinksFilter struct{}

func NewInternalLinksFilter() Filter {
	return internalLinksFilter{}
}

func (f internalLinksFilter) Use(message string) bool {
	return isIdyllURL(message)
}

func NewExternalLinksFilter() Filter {
	return externalLinksFilter{}
}

func (f externalLinksFilter) Use(message string) bool {
	return isContainsURL(message)
}
