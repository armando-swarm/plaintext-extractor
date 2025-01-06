package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	extractor := NewExtractor()
	tests := []struct {
		input    string
		expected string
	}{
		{`a<br>b`, "a\nb"},
		{`a<br><h1>b</h1>`, "a\nb\n"},
		{`<a href="https://example.com">link</a>`, "link"},
		{`<div>This is a <a href="https://example.com">link</a></div>`, "This is a link\n"},
		{"<div><h1>Heading 1</h1><h2>Heading 2</h2><ul><li>Item 1</li><li>Item 2</li></ul></div>", "Heading 1\nHeading 2\n- Item 1\n- Item 2\n"},
		{"<div><h1>Heading 1</h1><h2>Heading 2</h2><ol><li>Item 1</li><li>Item 2</li></ol></div>", "Heading 1\nHeading 2\n1. Item 1\n2. Item 2\n"},
		{"<p><span>a</span><span>b</span></p> c", "a b\nc"},
		{"a\n \nb", "a\nb"},
	}
	for _, test := range tests {
		output, err := extractor.PlainText(test.input)
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, test.expected, *output)
	}
}
