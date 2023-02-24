package filter

type internalLinksFilter struct{ baseFilter }

type externalLinksFilter struct{ baseFilter }

func NewInternalLinksFilter() Filter {
	return internalLinksFilter{
		baseFilter: baseFilter{tag: "nil", name: "no-internal-links"},
	}
}

func (f internalLinksFilter) Use(message string) bool {
	return isIdyllURL(message)
}

func NewExternalLinksFilter() Filter {
	return externalLinksFilter{
		baseFilter: baseFilter{tag: "nel", name: "no-external-links"},
	}
}

func (f externalLinksFilter) Use(message string) bool {
	return isContainsURL(message)
}
