package parser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type SummaryParser struct {
	raw       string
	lastLevel int
	lastPath  string
	inChapter bool
	content   string
	Urls      []string
	tokenizer *html.Tokenizer
}

// parse the raw html and return a error if something wrong
func (n *SummaryParser) Parse() error {
	reader := strings.NewReader(n.raw)
	n.tokenizer = html.NewTokenizer(reader)
	for {
		if n.tokenizer.Next() == html.ErrorToken {
			break
		}
		token := n.tokenizer.Token()
		if token.Type == html.StartTagToken && token.Data == "li" {
			n.content += n._parse(token, 0)
		}
	}
	return nil
}

// parse the list recursive
func (n *SummaryParser) _parse(token html.Token, level int) string {
	attrs := parseAttr(token.Attr)
	if class, ok := attrs["class"]; !ok || !strings.Contains(class, "chapter") {
		return ""
	}
	result := ""
	text := ""
	dataPath := attrs["data-path"]
	for {
		if n.tokenizer.Next() == html.ErrorToken {
			return result
		}
		token := n.tokenizer.Token()
		if token.Type == html.StartTagToken {
			if token.Data == "li" {
				result += n._parse(token, level+1)
			}
		} else if token.Type == html.TextToken {
			text += "" + strings.TrimSpace(token.Data)
		} else if token.Type == html.EndTagToken {
			if token.Data == "li" {
				result = fmt.Sprintf("%s*[%s](%s)\n", strings.Repeat("\t", level), text, dataPath) + result
				n.Urls = append(n.Urls, dataPath)
				return result
			}
		}
	}
}

// return the title of the page
func (n *SummaryParser) Title() string {
	return ""
}

// return the markdown format content of the summary
func (n *SummaryParser) Content() string {
	return n.content
}
