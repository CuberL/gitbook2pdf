package parser

import (
	"strings"

	"golang.org/x/net/html"
)

type Parser interface {
	Parse() error
	Title() string
	Content() string
}

// get the parser
func New(raw string, summary bool) Parser {
	if !summary {
		return &NormalParser{raw: raw}
	} else {
		return &SummaryParser{raw: raw, inChapter: false}
	}
}

func parseAttr(attrs []html.Attribute) map[string]string {
	result := map[string]string{}
	for _, attr := range attrs {
		result[attr.Key] = strings.Trim(attr.Val, " ")
	}
	return result
}
