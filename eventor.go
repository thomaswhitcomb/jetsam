package jetsam

import "log"

type eventor struct {
	out     chan<- string
	sources []string
	count   int
}

func (e *eventor) run() {
	for _, source := range e.sources {
		e.out <- source
		e.count++
	}
	// Tell any channel consumers that eventor is done.
	close(e.out)
	log.Printf("%d URLs queued for processors\n", len(e.sources))
}
