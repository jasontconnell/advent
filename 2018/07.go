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

	fmt.Println("Time", time.Since(startTime))
}

func sortList(list []*step) {
	for _, step := range list {
		sortList(step.prereqs)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].name < list[j].name })
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
