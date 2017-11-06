package parser

import (
	"strings"

	"golang.org/x/net/html"
)

var Printers map[string]Printer

func Init() {
	Printers = map[string]Printer{
		"h1": PrintTagH,
		"h2": PrintTagH,
		"h3": PrintTagH,
		"h4": PrintTagH,
		"h5": PrintTagH,
		"h6": PrintTagH,
		"li": PrintTagLi,
	}
}

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
