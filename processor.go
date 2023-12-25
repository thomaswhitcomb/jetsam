package jetsam

import (
	"bufio"
	"io"
	"log"
	"net/http"
)

type Loader func(url string) (io.Reader, error)
type DoneChanMsg struct {
	Text  string
	Count int
}
type processor struct {
	in     <-chan string
	out    chan<- string
	done   chan<- DoneChanMsg
	loader Loader
}

func (p *processor) init() {
	if p.loader == nil {
		p.loader = func(url string) (io.Reader, error) {
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Failed to GET %s\n", url)
				return nil, err
			}
			//defer resp.Body.Close()
			return resp.Body, nil
		}
	}
}

func (p *processor) run() {
	log.Printf("Processor starting\n")
	itemsLoaded := 0
	p.init()
	for {
		url, ok := <-p.in
		if !ok {
			log.Printf("Processor done. No URLs remaining.\n")
			break
		}
		ldr, err := p.loader(url)
		if err != nil {
			log.Printf("Failed to GET url: %s err: %v\n", url, err)
			continue
		}
		scanner := bufio.NewScanner(ldr)
		scanner.Scan()
		for scanner.Scan() {
			p.out <- scanner.Text()
			itemsLoaded++
		}
	}
	log.Printf("Processor terminating. Items loaded: %d\n", itemsLoaded)
	p.done <- DoneChanMsg{Text: "done", Count: itemsLoaded}
}
