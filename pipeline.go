package jetsam

import "log"

type Pipeline struct {
	eventer        eventer
	processor      []processor
	ProcessorCount int
	reducer        reducer
	sourceBacklog  chan string
	sourceLines    chan string
	processorDone  chan DoneChanMsg
	reducerDone    chan DoneChanMsg
	Sources        []string
	BufferSize     int
	Loader         Loader
	Reducer        ReducerFunc
}

func (ppl *Pipeline) Provision() {
	ppl.processor = []processor{}

	ppl.sourceLines = make(chan string, ppl.BufferSize)
	ppl.processorDone = make(chan DoneChanMsg, len(ppl.Sources))
	ppl.sourceBacklog = make(chan string, len(ppl.Sources))
	ppl.reducerDone = make(chan DoneChanMsg, 1)

	for i := 0; i < ppl.ProcessorCount; i++ {
		p := processor{
			in:     ppl.sourceBacklog,
			out:    ppl.sourceLines,
			done:   ppl.processorDone,
			loader: ppl.Loader,
		}
		ppl.processor = append(ppl.processor, p)
	}
	ppl.eventer = eventer{
		sources: ppl.Sources,
		out:     ppl.sourceBacklog,
	}
	ppl.reducer = reducer{
		in:   ppl.sourceLines,
		done: ppl.reducerDone,
		fn:   ppl.Reducer,
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
		totalLoaded = totalLoaded + msg.NumberProcessed
	}
	close(ppl.sourceLines)
	log.Printf("%d sourceLines loaded\n", totalLoaded)
	msg := <-ppl.reducerDone
	log.Printf("%d Items reduced\n", msg.NumberProcessed)
	log.Printf("%v\n", msg.Results)
}
