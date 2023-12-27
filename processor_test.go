package jetsam

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type MyReadCloser struct{}

func (MyReadCloser) Close() error {
	return nil
}
func (MyReadCloser) Read(p []byte) (n int, err error) {
	return 5, nil
}

func TestProcessorHappyPath(t *testing.T) {
	testLine := "this is a line of data"
	r := processor{
		in:   make(chan string, 1),
		out:  make(chan string, 1),
		done: make(chan DoneChanMsg, 1),
		loader: func(url string) (io.Reader, error) {
			rdr := bytes.NewReader([]byte("CSV header\n" + testLine))
			return rdr, nil
		},
		closer: MyReadCloser{},
	}

	r.in <- "fake url"
	close(r.in)

	r.run()
	resultLine := <-r.out
	msg := <-r.done

	if resultLine != testLine {
		t.Errorf("Unexpected line from chan: %s", resultLine)
	}
	if msg.NumberProcessed != 1 {
		t.Errorf("Incorrect NumberProcessed value: %d", msg.NumberProcessed)
	}
}
func TestProcessorFailedGET(t *testing.T) {
	r := processor{
		in:   make(chan string, 1),
		out:  make(chan string, 1),
		done: make(chan DoneChanMsg, 1),
		loader: func(url string) (io.Reader, error) {
			return nil, errors.New("test for bad URL")
		},
		closer: MyReadCloser{},
	}

	r.in <- "fake url"
	close(r.in)

	r.run()
	msg := <-r.done

	if msg.NumberProcessed != 0 {
		t.Errorf("Incorrect NumberProcessed value: %d", msg.NumberProcessed)
	}
}
