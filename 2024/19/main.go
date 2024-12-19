package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type state struct {
	pos   int
	sub   string
	total string
}

type statekey struct {
	pos int
	sub string
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	p, d := parse(in)
	return findValidDesigns(p, d)
}

func part2(in input) output {
	p, d := parse(in)
	for _, design := range d {
		countVariations(p, design)
	}
	return 0
}

func findValidDesigns(pmap map[string][]string, designs []string) int {
	min, max := minmaxlen(pmap)
	count := 0
	for _, design := range designs {
		if designValid(pmap, design, min, max) {
			count++
		}
	}
	return count
}

func minmaxlen(pmap map[string][]string) (int, int) {
	min, max := math.MaxInt32, math.MinInt32
	for s := range pmap {
		if len(s) < min {
			min = len(s)
		}
		if len(s) > max {
			max = len(s)
		}
	}
	return min, max
}

func countVariations(pmap map[string][]string, design string) {

}

func designValid(pmap map[string][]string, design string, min, max int) bool {
	valid := false
	queue := common.NewQueue[state, int]()
	initial := state{pos: 0, sub: "", total: ""}
	visit := make(map[statekey]bool)
	queue.Enqueue(initial)
	for queue.Any() {
		cur := queue.Dequeue()

		sk := statekey{pos: cur.pos, sub: cur.sub}
		if _, ok := visit[sk]; ok {
			continue
		}
		visit[sk] = true

		if cur.total == design {
			valid = true
			continue
		}

		mvs := getPossibles(design, cur, pmap, min, max)
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
	}
	return valid
}

func getPossibles(design string, cur state, pmap map[string][]string, min, max int) []state {
	list := []state{}
	for i := min; i <= max; i++ {
		if cur.pos+i > len(design) {
			continue
		}
		sub := design[cur.pos : cur.pos+i]
		if _, ok := pmap[sub]; ok {
			list = append(list, state{pos: cur.pos + i, sub: sub, total: cur.total + sub})
		}
	}
	return list
}

func parse(in []string) (map[string][]string, []string) {
	patterns := strings.Fields(strings.ReplaceAll(in[0], ",", ""))
	designs := in[2:]
	pmap := make(map[string][]string)
	for _, p := range patterns {
		pmap[p] = allSubstrings(p)
	}

	return pmap, designs
}

func allSubstrings(str string) []string {
	list := []string{}
	for i := 0; i < len(str); i++ {
		for j := 0; j < len(str); j++ {
			if i > j {
				continue
			}
			v := str[i : j+1]
			list = append(list, v)
		}
	}
	return list
}
