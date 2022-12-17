package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type valve struct {
	name     string
	flowrate int
	valves   []string
}

type release struct {
	minute, pressure, opened int
}

type max struct {
	pressure, flowrate, lastopen int
}

type state struct {
	valve    string
	prev     string
	minute   int
	released []release
	open     map[string]bool
	working  bool
}

type move struct {
	from, to string
	minutes  int
	open     bool
}

type diststate struct {
	steps int
	valve string
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 16 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	valves := parseInput(in)
	return solve("AA", 30, valves)
}

func part2(in input) output {
	return 0
}

func solve(start string, time int, valves map[string]valve) int {
	priority := getPriorityValves(valves)
	queue := common.NewQueue[state, int]()
	maxpressure := 0

	paths := calcPaths(valves)

	queue.Enqueue(state{valve: start, minute: 1, open: make(map[string]bool), released: []release{}, working: false})
	for queue.Any() {
		cur := queue.Dequeue()

		pressure := getReleased(cur.released, time)
		if pressure > maxpressure {
			maxpressure = pressure
		}

		if cur.working {
			cur.working = false
			cur.minute++
			queue.Enqueue(cur)
			continue
		}

		mvs := getMoves(cur, time-cur.minute, valves, paths, priority)
		for _, mv := range mvs {
			released := copySlice(cur.released)
			open := copyMap(cur.open)

			rel := release{minute: cur.minute + mv.minutes, pressure: valves[mv.to].flowrate}
			released = append(released, rel)
			open[mv.to] = true

			nstate := state{valve: mv.to, prev: mv.from, minute: cur.minute + mv.minutes, open: open, working: true, released: released}
			queue.Enqueue(nstate)
		}
	}

	return maxpressure
}

func getMoves(s state, timeleft int, valves map[string]valve, paths map[string]map[string]int, priority []string) []move {
	mvs := []move{}
	for _, r := range priority {
		if s.open[r] {
			continue
		}

		shortest := paths[s.valve][r]
		if shortest == math.MaxInt32 || shortest == 0 {
			continue
		}

		if shortest+1 > timeleft {
			continue
		}
		mv := move{from: s.valve, to: r, minutes: shortest, open: true}
		mvs = append(mvs, mv)
	}
	sort.Slice(mvs, func(i, j int) bool {
		return mvs[i].minutes < mvs[j].minutes
	})
	return mvs
}

func calcPaths(valves map[string]valve) map[string]map[string]int {
	paths := make(map[string]map[string]int)
	for k := range valves {
		paths[k] = make(map[string]int)
	}

	for pk := range paths {
		for pk2 := range paths {
			if pk == pk2 {
				continue
			}
			v := valves[pk2]
			if v.flowrate == 0 {
				continue
			}

			paths[pk][pk2] = shortestDistance(pk, pk2, valves)
		}
	}
	return paths
}

func getPriorityValves(valves map[string]valve) []string {
	list := []valve{}
	for _, v := range valves {
		if v.flowrate > 0 {
			list = append(list, v)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].flowrate > list[j].flowrate
	})

	s := []string{}
	for _, v := range list {
		s = append(s, v.name)
	}
	return s
}

func getReleased(list []release, minute int) int {
	total := 0
	for _, rel := range list {
		if rel.minute < minute {
			total += rel.pressure * (minute - rel.minute)
		}
	}
	return total
}

func getPressure(list []release, minute int) int {
	total := 0
	for _, rel := range list {
		if rel.minute <= minute {
			total += rel.pressure
		}
	}
	return total
}

func shortestDistance(from, to string, valves map[string]valve) int {
	queue := common.NewQueue[diststate, int]()
	queue.Enqueue(diststate{valve: from, steps: 0})
	visit := make(map[string]bool)
	shortest := math.MaxInt32

	for queue.Any() {
		cur := queue.Dequeue()

		if _, ok := visit[cur.valve]; ok {
			continue
		}
		visit[cur.valve] = true

		if cur.valve == to {
			if cur.steps < shortest {
				shortest = cur.steps
			}
			continue
		}

		if cur.steps >= shortest {
			continue
		}

		v := valves[cur.valve]
		for _, cv := range v.valves {
			st := diststate{steps: cur.steps + 1, valve: cv}
			queue.Enqueue(st)
		}
	}
	return shortest
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	nm := make(map[K]V)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

func keys[K comparable, V any](m map[K]V) []K {
	keys := []K{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func copySlice[T any](s []T) []T {
	cp := make([]T, len(s))
	copy(cp, s)
	return cp
}

func remove[T comparable](s []T, element T) []T {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == element {
			s = append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func parseInput(in input) map[string]valve {
	valves := []valve{}
	vmap := make(map[string]valve)
	reg := regexp.MustCompile("Valve ([A-Z]+) has flow rate=([0-9]+); tunnels? leads? to valves? ([A-Z, ]+)")
	for _, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) == 4 {
			rate, _ := strconv.Atoi(m[2])
			v := valve{name: m[1], flowrate: rate}
			sp := strings.Split(m[3], ", ")
			for _, s := range sp {
				v.valves = append(v.valves, s)
			}
			vmap[v.name] = v
			valves = append(valves, v)
		}
	}

	return vmap
}
