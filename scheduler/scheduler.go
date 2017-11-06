package scheduler

import (
	"fmt"
	"gitbook2pdf/parser"
	"io/ioutil"
	"net/http"
	"os"
)

type Scheduler struct {
	urls      chan (string)
	maxWorker int
	startUrl  string
}

func New(maxWorker int, startUrl string) *Scheduler {
	return &Scheduler{
		urls:      make(chan (string)),
		maxWorker: maxWorker,
		startUrl:  startUrl,
	}
}

func (s *Scheduler) Start() {
	// Get and save the summary
	summaryRaw, err := http.Get(s.startUrl)
	if err != nil {
		fmt.Printf("get summary failed: %s\n", err)
	}
	summaryHtml, err := ioutil.ReadAll(summaryRaw.Body)
	if err != nil {
		fmt.Printf("get summary failed: %s\n", err)
	}
	p := parser.New(string(summaryHtml), true)
	err = p.Parse()
	if err != nil {
		fmt.Printf("parse summary failed: %s\n", err)
	}
	err = ioutil.WriteFile("SUMMARY.md", []byte(p.Content()), os.ModeAppend)
	if err != nil {
		fmt.Printf("save summary failed: %s\n", err)
	}

	// Get the links
	urls := p.(*parser.SummaryParser).Urls
	for i := 0; i < s.maxWorker; i++ {
		go s.fetch()
	}
	for _, u := range urls {
		s.urls <- s.startUrl + "/" + u
	}
	close(s.urls)
}

func (s *Scheduler) fetch() {
	for url := range s.urls {
		if url == "" {
			return
		}
		contentRaw, err := http.Get(url)
		if err != nil {
			fmt.Printf("get content failed: %s\n", err.Error())
			continue
		}
		contentHtml, err := ioutil.ReadAll(contentRaw.Body)
		if err != nil {
			fmt.Printf("goroutine dead: %s\n", err.Error())
		}
		p := parser.New(string(contentHtml), false)
		err = p.Parse()
		if err != nil {
			fmt.Printf("parse error: %s\n", err.Error())
		}
		//		fmt.Println(p.Content())
	}
}
