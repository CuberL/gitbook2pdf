package parser

import (
	//	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type NormalParser struct {
	raw       string
	content   string
	tokenizer *html.Tokenizer
}

// parse the raw html and return a error when something wrong
func (n *NormalParser) Parse() error {
	n.tokenizer = html.NewTokenizer(strings.NewReader(n.raw))
	for {
		if n.tokenizer.Next() == html.ErrorToken {
			break
		}
		token := n.tokenizer.Token()
		if class, ok := parseAttr(token.Attr)["class"]; ok {
			if strings.Contains(class, "markdown-section") {
				n.content = n._parse(token)
			}
		}
	}
	return nil
}

// return the title of the page
func (n *NormalParser) Title() string {
	return ""
}

// return the markdown format content of the page
func (n *NormalParser) Content() string {
	return n.content
}

func (n *NormalParser) _parse(token html.Token) string {
	mdData := ""
	mdTag := ""
	newLine := false
	switch token.Data {
	case "li":
		mdTag = "*"
		newLine = true
	case "h1":
		mdTag = "#"
		newLine = true
	case "h2":
		mdTag = "##"
		newLine = true
	case "h3":
		mdTag = "###"
		newLine = true
	case "h4":
		mdTag = "####"
		newLine = true
	case "h5":
		mdTag = "#####"
		newLine = true
	case "h6":
		mdTag = "######"
		newLine = true
	case "code":
		mdTag = "```"
	case "pre":
		newLine = true
	default:
		mdTag = ""
	}
	for {
		n.tokenizer.Next()
		nextToken := n.tokenizer.Token()
		switch nextToken.Type {
		case html.TextToken:
			mdData += nextToken.Data
		case html.StartTagToken:
			mdData += n._parse(nextToken)
		case html.EndTagToken:
			line := ""
			if newLine {
				line = "\n"
			} else {
				line = ""
			}
			return mdTag + " " + strings.TrimSpace(mdData) + " " + mdTag + line
		}
	}

}
