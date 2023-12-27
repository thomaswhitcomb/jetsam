package jetsam

import (
	"testing"
)

func TestEventor(t *testing.T) {
	s := []string{"a", "b", "c"}
	e := eventor{
		out:     make(chan string, len(s)),
		sources: s,
	}
	e.run()
	if e.count != len(s) {
		t.Errorf("Bad resulting length: %d\n", e.count)
	}
	if len(e.out) != len(s) {
		t.Errorf("Bad resulting channel length: %d\n", len(e.out))
	}
}
