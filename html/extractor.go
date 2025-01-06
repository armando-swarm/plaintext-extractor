package html

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Extractor represents an HTML-specific plain text extractor.
type Extractor struct {
	blockTags map[string]bool
}

// NewExtractor creates a new HTMLExtractor instance.
func NewExtractor(otherBlockTags ...string) *Extractor {
	uniqueBlockTags := map[string]bool{}
	for _, tag := range blockTags {
		uniqueBlockTags[tag] = true
	}
	for _, tag := range otherBlockTags {
		uniqueBlockTags[tag] = true
	}

	return &Extractor{blockTags: uniqueBlockTags}
}

// PlainText extracts plain text from the input HTML string.
func (e *Extractor) PlainText(input string) (*string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return nil, err
	}

	var plainText strings.Builder
	e.extractText(&plainText, doc, 0)

	output := plainText.String()
	len := plainText.Len()
	var i int
	for i = 0; i < len; i++ {
		if !e.isSpace(output[i]) {
			break
		}
	}

	output = output[i:]
	output = string(regexp.MustCompile("\\s*\n+\\s*").ReplaceAll([]byte(output), []byte("\n")))
	return &output, nil
}

func (e *Extractor) isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// Recursively extract plain text from the HTML nodes.
func (e *Extractor) extractText(plainText *strings.Builder, node *html.Node, idx int) {
	liType := e.listItemType(node)
	if liType == OrderedListItem {
		plainText.WriteString(fmt.Sprintf("%d. ", idx))
	} else if liType == UnorderedListItem {
		plainText.WriteString("- ")
	}

	if node.Type == html.TextNode {
		if strings.TrimSpace(node.Data) != "" {
			plainText.WriteString(node.Data)
		}
	}
	if node.DataAtom.String() == "br" {
		plainText.WriteString("\n")
		return
	}

	i := 0
	var isList bool = node.DataAtom.String() == "ul" || node.DataAtom.String() == "ol"
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isList {
			i++
		}
		e.extractText(plainText, child, i)
	}
	if found := e.blockTags[node.DataAtom.String()]; found {
		plainText.WriteString("\n")
	}
}

type ListItemType int

const (
	Unknown           ListItemType = iota
	UnorderedListItem ListItemType = 1
	OrderedListItem   ListItemType = 2
)

func (e *Extractor) listItemType(node *html.Node) ListItemType {
	if node.DataAtom.String() != "li" {
		return Unknown
	}

	for p := node.Parent; p != nil; p = p.Parent {
		if p.DataAtom.String() == "ul" {
			return UnorderedListItem
		}
		if p.DataAtom.String() == "ol" {
			return OrderedListItem
		}
	}

	return Unknown
}
