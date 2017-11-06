package scheduler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuberl/gitbook2pdf/parser"
)

type Scheduler struct {
	paths     chan (string)
	maxWorker int
	startUrl  string
	storeDir  string
}

func New(maxWorker int, startUrl string, storeDir string) *Scheduler {
	return &Scheduler{
		paths:     make(chan (string)),
		maxWorker: maxWorker,
		startUrl:  startUrl,
		storeDir:  storeDir,
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
	p := parser.NewSummaryParser(string(summaryHtml))
	err = p.Parse()
	if err != nil {
		fmt.Printf("parse summary failed: %s\n", err)
	}
	err = ioutil.WriteFile("SUMMARY.md", []byte(p.Content()), 0666)
	if err != nil {
		fmt.Printf("save summary failed: %s\n", err)
	}

	// Get the links
	urls := p.Urls
	for i := 0; i < s.maxWorker; i++ {
		go s.fetch()
	}
	for _, u := range urls {
		s.paths <- u
	}
	close(s.paths)
}

func (s *Scheduler) fetch() {
	for path := range s.paths {
		if path == "" {
			return
		}
		url := s.startUrl + "/" + path
		contentRaw, err := http.Get(url)
		if err != nil {
			fmt.Printf("get content failed: %s\n", err.Error())
			continue
		}
		contentHtml, err := ioutil.ReadAll(contentRaw.Body)
		if err != nil {
			fmt.Printf("goroutine dead: %s\n", err.Error())
		}
		p := parser.NewParser(string(contentHtml))
		err = p.Parse()
		if err != nil {
			fmt.Printf("parse error: %s\n", err.Error())
		}
		path = s.storeDir + "/" + path
		baseDir := filepath.Dir(path)
		fmt.Println(baseDir)
		_, err = os.Stat(baseDir)
		if err != nil {
			err = os.MkdirAll(baseDir, 0777)
			if err != nil {
				fmt.Println("mkdir failed.")
				continue
			}
		}
		ioutil.WriteFile(strings.Replace(path, ".html", ".md", 1), []byte(p.Content()), 0666)
	}
}
