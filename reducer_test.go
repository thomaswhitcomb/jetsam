package jetsam

import (
	"testing"
)

func TestReducer(t *testing.T) {
	itRan := false

	r := reducer{}
	r.fn = func(in <-chan string, out chan<- DoneChanMsg) {
		itRan = true
	}
	r.run()

	if !itRan {
		t.Errorf("Reducer function not called")
	}
}
