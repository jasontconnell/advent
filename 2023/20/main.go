package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type pulse int

const (
	nopulse pulse = iota
	low
	high
)

func (p pulse) String() string {
	s := "none"
	if p == low {
		s = "low"
	} else if p == high {
		s = "high"
	}
	return s
}

type state int

const (
	nostate state = iota
	off
	on
)

func (st state) String() string {
	s := "nostate"
	if st == off {
		s = "off"
	} else if st == on {
		s = "on"
	}
	return s
}

type module struct {
	name        string
	flipflop    bool
	conjunction bool
	targets     []string
	watches     []string
	received    map[string]pulse
	state       state
	output      bool
}

type pulseState struct {
	name  string
	from  string
	pulse pulse
}

func (m *module) String() string {
	return fmt.Sprintf("%s: flipflop: %t conj: %t state: %v", m.name, m.flipflop, m.conjunction, m.state)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	modules := parseInput(in)
	low, high := pushButton(modules, "broadcaster", low, 1000)
	return low * high
}

func part2(in input) output {
	return 0
}

func pushButton(mods map[string]*module, startname string, startpulse pulse, times int) (int, int) {
	l, h := 0, 0
	for i := 0; i < times; i++ {
		queue := common.NewQueue[pulseState, int]()
		queue.Enqueue(pulseState{name: startname, pulse: startpulse})
		for queue.Any() {
			cur := queue.Dequeue()

			if cur.pulse == low {
				l++
			} else if cur.pulse == high {
				h++
			}

			m, _ := mods[cur.name]
			if m == nil {
				continue
			}
			res := queuePulses(mods, m, cur.from, cur.pulse, cur.name == startname)
			for _, r := range res {

				queue.Enqueue(r)
			}
		}
	}
	return l, h
}

func queuePulses(mm map[string]*module, m *module, from string, p pulse, isorigin bool) []pulseState {
	queue := []pulseState{}
	if m.flipflop {
		np := nopulse
		if p == low && m.state == off {
			np = high
			m.state = on
		} else if p == low && m.state == on {
			np = low
			m.state = off
		}

		if np != nopulse {
			for _, tg := range m.targets {
				queue = append(queue, pulseState{name: tg, from: m.name, pulse: np})
			}
		}
	} else if m.conjunction {
		if p != nopulse {
			m.received[from] = p
		}

		np := nopulse
		allhigh := true
		for _, w := range m.watches {
			if p, ok := m.received[w]; ok {
				np = p
				if np != high {
					allhigh = false
				}
			}
		}

		np = high
		if allhigh {
			np = low
		}

		for _, tg := range m.targets {
			queue = append(queue, pulseState{name: tg, pulse: np, from: m.name})
		}
	} else if isorigin {
		np := p
		for _, tg := range m.targets {
			queue = append(queue, pulseState{name: tg, pulse: np, from: "origin"})
		}
	}
	return queue
}

func parseInput(in input) map[string]*module {
	ms := []*module{}
	hasoutput := false
	for _, line := range in {
		if !hasoutput {
			hasoutput = strings.Index(line, "output") != -1
		}
		sp := strings.Split(line, " -> ")
		mname := strings.Trim(sp[0], " ")
		flipflop := mname[0] == '%'
		conjunction := mname[0] == '&'

		if conjunction || flipflop {
			mname = mname[1:]
		}

		mod := &module{name: mname, flipflop: flipflop, conjunction: conjunction}
		if flipflop {
			mod.state = off
		}
		tgs := strings.Split(sp[1], ", ")
		mod.targets = tgs
		ms = append(ms, mod)
	}

	if hasoutput {
		ms = append(ms, &module{name: "output", output: true})
	}

	mm := make(map[string]*module)
	for _, m := range ms {
		mm[m.name] = m
	}

	for _, m := range ms {
		if m.conjunction {
			m.received = make(map[string]pulse)
			for _, ss := range ms {
				for _, t := range ss.targets {
					if t == m.name {
						m.watches = append(m.watches, ss.name)
						m.received[ss.name] = low
					}
				}
			}
		}
	}

	return mm
}
