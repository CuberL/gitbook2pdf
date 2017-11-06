package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func NewSummaryParser(raw string) *SummaryParser {
	return &SummaryParser{raw: raw, inChapter: false}
}

func NewParser(raw string) *NormalParser {
	return &NormalParser{raw: raw}
}

func parseAttr(attrs []html.Attribute) map[string]string {
	result := map[string]string{}
	for _, attr := range attrs {
		result[attr.Key] = strings.Trim(attr.Val, " ")
	}
	return result
}
