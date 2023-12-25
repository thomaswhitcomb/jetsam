package jetsam

import "log"

type eventer struct {
	out     chan<- string
	sources []string
	count   int
}

func (e *eventer) run() {
	for _, source := range e.sources {
		e.out <- source
		e.count++
	}
	close(e.out)
	log.Printf("%d URLs queued for processors\n", len(e.sources))
}
