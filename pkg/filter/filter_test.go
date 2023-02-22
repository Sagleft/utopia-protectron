package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelsFilter(t *testing.T) {
	f := NewChannelsFilter()

	assert.False(t, f.Use("test"))
	assert.False(t, f.Use("afff"))
	assert.True(t, f.Use("af530dc7ba83c04bf7d3a02c5d8a9cf3"))
	assert.False(t, f.Use("zf530dc7ba83c04bf7d3a02c5d8a9cf3"))
}

func TestExternalLinksFilter(t *testing.T) {
	f := NewExternalLinksFilter()

	assert.False(t, f.Use("example"))
	assert.False(t, f.Use("utopia://example"))
	assert.True(t, f.Use("http://example.com"))
	assert.True(t, f.Use("https://example.com"))
	assert.True(t, f.Use("link: http://example.com"))
	assert.True(t, f.Use("link: https://example.com"))
}

func TestInternalLinksFilter(t *testing.T) {
	f := NewInternalLinksFilter()

	assert.False(t, f.Use("example"))
	assert.False(t, f.Use("http://example.com"))
	assert.False(t, f.Use("https://example.com"))
	assert.False(t, f.Use("link: https://example.com"))
	assert.True(t, f.Use("link: utopia://idyll"))
	assert.True(t, f.Use("utopia://example"))
}
