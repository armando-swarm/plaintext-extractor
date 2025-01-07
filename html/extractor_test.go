package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const rawHTML = `<p>This is a <b>bold</b> text inside a paragraph.</p>
<p>This is an <i>italicized</i> text inside another paragraph.</p>
<p>
    This is a <b><i>bold and italicized</i></b> text inside a paragraph.
</p>
<p>
    Nested <i> and maybe italic <b>and bold</b> and not</i>.
</p>

<ul>
    <li>This is a <b>bold</b> list item.</li>
    <li>This is an <i>italicized</i> list item.</li>
    <li>
        This list item contains both <b><i>bold and italicized</i></b> text.
    </li>
    <li>
        Another example with nested tags: <b>This is <i>bold and italic</i> text</b>.
    </li>
	<li>
	    This is an example of a list item
		with a linebreak
	</li>
	<li>
	     And what happens if another item <b>follows</b> the linebreak
	</li>
</ul>

<p>
  Testing <i>hyperlinks</i> <a href="https://www.google.com">click <b>me <i>quickly!</i> or else</b> go boom</a> inside a <b>paragraph with<i> complex </i>formatting</b>.
</p>

<p>
  <b>Testing hyperlinks <a href="https://www.swarm.work"><i>inheritance</i> stuff</a> that inherit outer tags</b>
</p>
`

const plainTextHTML = `This is a bold text inside a paragraph.
This is an italicized text inside another paragraph.
This is a bold and italicized text inside a paragraph.
Nested  and maybe italic and bold and not.
- This is a bold list item.
- This is an italicized list item.
- This list item contains both bold and italicized text.
- Another example with nested tags: This is bold and italic text.
- This is an example of a list item
with a linebreak
- And what happens if another item follows the linebreak
Testing hyperlinksclick me quickly! or else go boom inside a paragraph with complex formatting.
Testing hyperlinks inheritance stuff that inherit outer tags
`

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
		{"<p><span>a</span><span>b</span></p> c", "ab\nc"},
		{"a\n \nb", "a\nb"},
		{"Nested <i>and maybe italic <b>and bold</b> and not</i>.", "Nested and maybe italic and bold and not."},
		{rawHTML, plainTextHTML},
	}
	for _, test := range tests {
		output, err := extractor.PlainText(test.input)
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, test.expected, *output)
	}
}
