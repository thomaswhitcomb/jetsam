package jetsam

import "log"

type Pipeline struct {
	eventer        eventer
	processor      []processor
	ProcessorCount int
	reducer        reducer
	work           chan string
	items          chan string
	processorDone  chan doneChanMsg
	reducerDone    chan doneChanMsg
	Sources        []string
	ItemDepth      int
	Loader         Loader
}

func (ppl *Pipeline) Provision() {
	ppl.processor = []processor{}

	ppl.items = make(chan string, ppl.ItemDepth)
	ppl.processorDone = make(chan doneChanMsg, len(ppl.Sources))
	ppl.work = make(chan string, len(ppl.Sources))
	ppl.reducerDone = make(chan doneChanMsg, 1)

	for i := 0; i < ppl.ProcessorCount; i++ {
		p := processor{
			in:     ppl.work,
			out:    ppl.items,
			done:   ppl.processorDone,
			loader: ppl.Loader,
		}
		ppl.processor = append(ppl.processor, p)
	}
	ppl.eventer = eventer{
		sources: ppl.Sources,
		out:     ppl.work,
	}
	ppl.reducer = reducer{
		in:   ppl.items,
		done: ppl.reducerDone,
	}
}
func (ppl *Pipeline) Run() {
	go ppl.reducer.run()
	ppl.eventer.run()
	for _, task := range ppl.processor {
		task := task
		go task.run()
	}
	totalLoaded := 0
	for i := 0; i < len(ppl.processor); i++ {
		msg := <-ppl.processorDone
		totalLoaded = totalLoaded + msg.count
	}
	close(ppl.items)
	log.Printf("%d items loaded\n", totalLoaded)
	msg := <-ppl.reducerDone
	log.Printf("%d Items reduced\n", msg.count)
	log.Printf("%s\n", msg.text)
}
