// html2md project main.go
package main

import (
	"gitbook2pdf/scheduler"
)

// The Application has these three parts:
// Fetcher -- Fetch pages from a start link
// Parser -- Convert a HTML page to Markdown file
// Scheduler
func main() {
	s := scheduler.New(1, "http://localhost:4000")
	s.Start()
}
