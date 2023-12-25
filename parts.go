package jetsam

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

/*
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
*/

type Loader func(url string) (io.Reader, error)
type doneChanMsg struct {
	text  string
	count int
}
type processor struct {
	in     <-chan string
	out    chan<- string
	done   chan<- doneChanMsg
	loader Loader
}

func (p *processor) init() {
	if p.loader == nil {
		p.loader = func(url string) (io.Reader, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			return resp.Body, nil
		}
	}
}

func (p *processor) run() {
	log.Printf("Processor starting\n")
	itemsLoaded := 0
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
		scanner := bufio.NewScanner(ldr)
		scanner.Scan()
		for scanner.Scan() {
			p.out <- scanner.Text()
			itemsLoaded++
		}
	}
	log.Printf("Processor terminating. Items loaded: %d\n", itemsLoaded)
	p.done <- doneChanMsg{text: "done", count: itemsLoaded}
}

type reducer struct {
	in   <-chan string
	done chan<- doneChanMsg
}
type ageName struct {
	fname string
	lname string
}
type ageNameMap map[int][]ageName

func findMedium(ages ageNameMap) (float64, []ageName) {
	agesonly := []int{}
	n := 0
	for k, _ := range ages {
		agesonly = append(agesonly, k)
		n++
	}
	sort.Ints(agesonly)
	if n%2 == 0 {
		slot := int(math.Floor(float64(n)/2) - 1)

		avg := float64(agesonly[slot]+agesonly[slot+1]) / 2
		return avg, nil
	} else {
		slot := int(math.Floor(float64(n) / 2.0))
		return float64(agesonly[slot]), ages[agesonly[slot]]
	}
}
func (r *reducer) run() {
	n := 0
	average := 0
	ages := map[int][]ageName{}
	for {
		line, ok := <-r.in
		if ok {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				log.Printf("Less than 3: %s\n", line)
				continue
			}
			age, err := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err != nil {
				log.Printf("Age conversion failed. |%s| %v\n", strings.TrimSpace(parts[2]), err)
				continue
			}
			an := ageName{
				fname: strings.TrimSpace(parts[0]),
				lname: strings.TrimSpace(parts[1]),
			}
			if val, ok := ages[age]; ok {
				ages[age] = append(val, an)
			} else {
				ages[age] = []ageName{an}
			}
			average = average + age

			n++
		} else {
			fmt.Printf("Reducing input queue empty. Stopping\n")
			break
		}
	}
	age, who := findMedium(ages)
	if who != nil {
		r.done <- doneChanMsg{
			text: fmt.Sprintf(
				"Average: %f | Medium:%f for %s %s.",
				float64(average)/float64(n),
				age,
				who[0].fname,
				who[0].lname),
			count: n,
		}
	} else {
		r.done <- doneChanMsg{
			text: fmt.Sprintf(
				"Average: %f | Medium:%f ",
				float64(average)/float64(n),
				age),
			count: n,
		}
	}
}
