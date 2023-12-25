package jetsam

type ReducerFunc func(in <-chan string, don chan<- DoneChanMsg)
type reducer struct {
	in   <-chan string
	done chan<- DoneChanMsg
	fn   ReducerFunc
}

func (r *reducer) run() {
	r.fn(r.in, r.done)
}
