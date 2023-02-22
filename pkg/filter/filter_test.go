package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoLinksFilter(t *testing.T) {
	f := NewLinksFilter()

	assert.False(t, f.Use("test"))
	assert.False(t, f.Use("afff"))
	assert.True(t, f.Use("af530dc7ba83c04bf7d3a02c5d8a9cf3"))
	assert.False(t, f.Use("zf530dc7ba83c04bf7d3a02c5d8a9cf3"))
}
