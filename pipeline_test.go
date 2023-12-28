package jetsam

import (
	"bytes"
	"io"
	"testing"
)

func TestPipelineProvision(t *testing.T) {
	ppl := Pipeline{
		ProcessorCount: 2,
		BufferSize:     100,
		Sources:        []string{"1", "2", "3"},
		Loader: func(url string) (io.Reader, error) {
			return nil, nil
		},
		Reducer: func(in <-chan string, out chan<- DoneChanMsg) {},
	}
	ppl.Provision()

	if len(ppl.processor) != 2 {
		t.Errorf("Incorrect processor count: %d", len(ppl.processor))
	}
}

// 1 processor, 1 URL and 1 reduced lines
func TestPipelineRunHappy1(t *testing.T) {
	testLine := "this is a line of data"
	ppl := Pipeline{
		ProcessorCount: 1,
		BufferSize:     100,
		Sources:        []string{"a url"},
		Loader: func(url string) (io.Reader, error) {
			rdr := bytes.NewReader([]byte("CSV header\n" + testLine))
			return rdr, nil
		},
		Reducer: func(in <-chan string, done chan<- DoneChanMsg) {
			n := 0
			for {
				_, ok := <-in
				if ok {
					n++
				} else {
					break
				}
			}
			done <- DoneChanMsg{
				Results:         DoneParms{},
				NumberProcessed: n,
			}
		},
	}
	ppl.Provision()
	ppl.processor[0].closer = MyReadCloser{}
	msg := ppl.Run()
	if msg.NumberProcessed != 1 {
		t.Errorf("Incorrect NumberProcessed: %d", msg.NumberProcessed)
	}
}

// 3 processor, 1 URL and 1 reduced lines
func TestPipelineRunHappy2(t *testing.T) {
	testLine := "this is a line of data"
	ppl := Pipeline{
		ProcessorCount: 3,
		BufferSize:     100,
		Sources:        []string{"a url"},
		Loader: func(url string) (io.Reader, error) {
			rdr := bytes.NewReader([]byte("CSV header\n" + testLine))
			return rdr, nil
		},
		Reducer: func(in <-chan string, done chan<- DoneChanMsg) {
			n := 0
			for {
				_, ok := <-in
				if ok {
					n++
				} else {
					break
				}
			}
			done <- DoneChanMsg{
				Results:         DoneParms{},
				NumberProcessed: n,
			}
		},
	}
	ppl.Provision()
	for i := 0; i < len(ppl.processor); i++ {
		ppl.processor[i].closer = MyReadCloser{}
	}
	msg := ppl.Run()
	if msg.NumberProcessed != 1 {
		t.Errorf("Incorrect NumberProcessed: %d", msg.NumberProcessed)
	}
}

// 3 processor, 1 URL and 3 reduced lines
func TestPipelineRunHappy3(t *testing.T) {
	testLine := "this is a line of data"
	ppl := Pipeline{
		ProcessorCount: 3,
		BufferSize:     100,
		Sources:        []string{"a url"},
		Loader: func(url string) (io.Reader, error) {
			rdr := bytes.NewReader(
				[]byte("CSV header\n" +
					testLine + "\n" +
					testLine + "\n" + testLine))
			return rdr, nil
		},
		Reducer: func(in <-chan string, done chan<- DoneChanMsg) {
			n := 0
			for {
				_, ok := <-in
				if ok {
					n++
				} else {
					break
				}
			}
			done <- DoneChanMsg{
				Results:         DoneParms{},
				NumberProcessed: n,
			}
		},
	}
	ppl.Provision()
	for i := 0; i < len(ppl.processor); i++ {
		ppl.processor[i].closer = MyReadCloser{}
	}
	msg := ppl.Run()
	if msg.NumberProcessed != 3 {
		t.Errorf("Incorrect NumberProcessed: %d", msg.NumberProcessed)
	}
}

// 3 processor, 3 URL and 3 reduced lines
func TestPipelineRunHappy4(t *testing.T) {
	testLine := "this is a line of data"
	ppl := Pipeline{
		ProcessorCount: 3,
		BufferSize:     100,
		Sources:        []string{"url1", "url2", "url3"},
		Loader: func(url string) (io.Reader, error) {
			rdr := bytes.NewReader(
				[]byte("CSV header\n" +
					testLine + "\n" +
					testLine + "\n" + testLine))
			return rdr, nil
		},
		Reducer: func(in <-chan string, done chan<- DoneChanMsg) {
			n := 0
			for {
				_, ok := <-in
				if ok {
					n++
				} else {
					break
				}
			}
			done <- DoneChanMsg{
				Results:         DoneParms{},
				NumberProcessed: n,
			}
		},
	}
	ppl.Provision()
	for i := 0; i < len(ppl.processor); i++ {
		ppl.processor[i].closer = MyReadCloser{}
	}
	msg := ppl.Run()
	if msg.NumberProcessed != 9 {
		t.Errorf("Incorrect NumberProcessed: %d", msg.NumberProcessed)
	}
}

// 1 processor, 3 URL and 3 reduced lines
func TestPipelineRunHappy5(t *testing.T) {
	testLine := "this is a line of data"
	ppl := Pipeline{
		ProcessorCount: 1,
		BufferSize:     100,
		Sources:        []string{"url1", "url2", "url3"},
		Loader: func(url string) (io.Reader, error) {
			rdr := bytes.NewReader(
				[]byte("CSV header\n" +
					testLine + "\n" +
					testLine + "\n" + testLine))
			return rdr, nil
		},
		Reducer: func(in <-chan string, done chan<- DoneChanMsg) {
			n := 0
			for {
				_, ok := <-in
				if ok {
					n++
				} else {
					break
				}
			}
			done <- DoneChanMsg{
				Results:         DoneParms{},
				NumberProcessed: n,
			}
		},
	}
	ppl.Provision()
	for i := 0; i < len(ppl.processor); i++ {
		ppl.processor[i].closer = MyReadCloser{}
	}
	msg := ppl.Run()
	if msg.NumberProcessed != 9 {
		t.Errorf("Incorrect NumberProcessed: %d", msg.NumberProcessed)
	}
}
