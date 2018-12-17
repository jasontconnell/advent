package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"
	//"strconv"
	//"strings"
	//"math"
)

var input = "07.txt"

type step struct {
	name      string
	prereqs   []*step
	pnames    []string
	completed bool
	inprocess bool
}

func (s *step) String() string {
	return s.name
}

type worker struct {
	step *step
	busy bool
	left int
}

func (s *step) Completed() bool {
	subcomplete := true

	if len(s.prereqs) == 0 {
		return true
	}
	for _, sub := range s.prereqs {
		subcomplete = subcomplete && sub.Completed()
	}
	return s.completed && subcomplete
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	m := make(map[string]*step)
	for scanner.Scan() {
		var txt = scanner.Text()
		getStep(m, txt)
	}

	list := []*step{}
	for _, v := range m {
		list = append(list, v)
	}

	list = mapPrereqs(list)
	sortList(list)

	order := process(list)
	fmt.Println("Process order: " + order)

	resetList(list)
	duration := processParallel(list, 5, 60)
	fmt.Println("Process time for jobs:", duration)

	fmt.Println("Time", time.Since(startTime))
}

func sortList(list []*step) {
	for _, step := range list {
		sortList(step.prereqs)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].name < list[j].name })
}

func resetList(list []*step) {
	for _, step := range list {
		step.completed = false
		step.inprocess = false
		resetList(step.prereqs)
	}
}

func process(steps []*step) string {
	p := ""
	n := getNext(steps)
	for n != nil {
		p += n.name
		n.completed = true
		n = getNext(steps)
	}
	return p
}

func processParallel(steps []*step, numWorkers, duration int) int {
	workers := []*worker{}
	for i := 0; i < numWorkers; i++ {
		w := &worker{busy: false}
		workers = append(workers, w)
	}
	i := 0
	complete := false
	for !complete {
		loop(steps, workers, duration)
		c := unfinished(steps)
		complete = c == 0
		if !complete {
			i++
		}
	}
	return i
}

func unfinished(steps []*step) int {
	i := len(steps)
	for _, s := range steps {
		if s.completed {
			i--
		}
	}
	return i
}

func loop(steps []*step, workers []*worker, duration int) {
	for i := 0; i < len(workers); i++ {
		w := workers[i]
		if !w.busy {
			workerReady(w, steps, duration)
		} else {
			w.left--
			if w.left == 0 {
				w.step.completed = true
				w.step.inprocess = false
				workerReady(w, steps, duration)
			}
		}
	}
}

func workerReady(w *worker, steps []*step, duration int) {
	s := nextToProcess(steps)
	if s != nil {
		ttl := duration + 1 + int(s.name[0]) - int('A')
		w.step = s
		w.step.inprocess = true
		w.busy = true
		w.left = ttl
	}
}

func nextToProcess(steps []*step) *step {
	if len(steps) == 0 {
		return nil
	}

	ready := []*step{}
	for _, s := range steps {
		sub := nextToProcess(s.prereqs)
		if sub != nil {
			ready = append(ready, sub)
		} else if !s.completed && !s.inprocess && unfinished(s.prereqs) == 0 {
			ready = append(ready, s)
		}
	}

	if len(ready) > 0 {
		sortList(ready)
		return ready[0]
	}

	return nil
}

func getNext(steps []*step) *step {
	ready := []*step{}
	for _, s := range steps {
		sub := getNext(s.prereqs)
		if sub != nil {
			ready = append(ready, sub)
		} else if !s.completed {
			ready = append(ready, s)
		}
	}

	if len(ready) > 0 {
		sortList(ready)
		return ready[0]
	}

	return nil
}

func getStep(m map[string]*step, line string) {
	reg := regexp.MustCompile("^Step ([A-Z]) must be finished before step ([A-Z]) can begin.$")
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		nm := groups[1]
		nm2 := groups[2]
		if ms, ok := m[nm]; ok {
			ms.pnames = append(ms.pnames, nm2)
		} else {
			s := &step{name: nm, pnames: []string{nm2}, prereqs: []*step{}}
			m[nm] = s
		}
	}
}

func mapPrereqs(list []*step) []*step {
	m := make(map[string]*step)
	for _, s := range list {
		m[s.name] = s
	}

	for _, s := range list {
		for _, pr := range s.pnames {
			if s1, ok := m[pr]; ok {
				s1.prereqs = append(s1.prereqs, s)
			} else {
				nodef := &step{name: pr, completed: false, prereqs: []*step{s}}
				m[pr] = nodef
				list = append(list, nodef)
			}
		}
	}

	return list
}
