package parser

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Printer func(string, html.Token) string

func PrintTagLi(data string, token html.Token) string {
	return "* " + data + "\n"
}

func PrintTagH(data string, token html.Token) string {
	level, err := strconv.Atoi(string(token.Data[1]))
	if err != nil {
		return ""
	}
	return strings.Repeat("#", level) + " " + data
}
