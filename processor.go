package jetsam

import (
	"bufio"
	"io"
	"log"
	"net/http"
)

type Loader func(url string) (io.Reader, error)
type DoneParms map[string]string
type DoneChanMsg struct {
	Results         DoneParms
	NumberProcessed int
}
type processor struct {
	in     chan string
	out    chan string
	done   chan DoneChanMsg
	loader Loader
	closer io.ReadCloser
}

func (p *processor) init() {
	if p.loader == nil {
		p.loader = func(url string) (io.Reader, error) {
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Failed to GET %s\n", url)
				return nil, err
			}
			p.closer = resp.Body
			return resp.Body, nil
		}
	}
}

func (p *processor) run() {
	log.Printf("Processor starting\n")
	linesLoaded := 0
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
		defer p.closer.Close()
		scanner := bufio.NewScanner(ldr)
		scanner.Scan() // skip csv header
		for scanner.Scan() {
			p.out <- scanner.Text()
			linesLoaded++
		}
	}
	log.Printf("Processor terminating. Lines loaded: %d\n", linesLoaded)
	p.done <- DoneChanMsg{
		Results:         DoneParms{},
		NumberProcessed: linesLoaded,
	}
}
