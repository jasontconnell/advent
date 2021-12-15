package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int64

type mapping struct {
	pair, insert string
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
	fmt.Fprintln(w, "--2021 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	w, rules := parseInput(in)
	n := runPolymer(w, rules, 10)
	min, max := minmax(n)
	return max - min
}

func part2(in input) output {
	w, rules := parseInput(in)
	n := runPolymer(w, rules, 40)
	min, max := minmax(n)
	return max - min
}

func minmax(m map[byte]int64) (int64, int64) {
	min, max := int64(math.MaxInt64), int64(math.MinInt64)

	for _, v := range m {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func runPolymer(w string, rules []mapping, steps int) map[byte]int64 {
	m := make(map[byte]int64)
	for x := 0; x < len(w); x++ {
		m[w[x]]++
	}

	pairs := make(map[string]int64)
	for i := 0; i < len(w)-1; i++ {
		p := string(w[i : i+2])
		pairs[p]++
	}

	rulemap := make(map[string]string)
	for _, r := range rules {
		rulemap[r.pair] = r.insert
	}

	for i := 0; i < steps; i++ {
		ns := map[string]int64{}
		for pp, v := range pairs {
			r, ok := rulemap[pp]
			if !ok {
				continue
			}

			upd := []string{string(pp[0]) + string(r[0]), string(r[0]) + string(pp[1])}
			for _, u := range upd {
				ns[u] += v
			}
			m[r[0]] += v
		}

		pairs = make(map[string]int64)
		for k, v := range ns {
			pairs[k] = v
		}
	}
	return m
}

func parseInput(in input) (string, []mapping) {
	wmode := true
	mappings := []mapping{}
	word := ""
	for _, line := range in {
		if line == "" {
			wmode = false
			continue
		}

		if wmode {
			word = line
		} else {
			p := strings.Split(line, "->")
			m := mapping{pair: strings.Trim(p[0], " "), insert: strings.Trim(p[1], " ")}
			mappings = append(mappings, m)
		}
	}
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].pair < mappings[j].pair
	})
	return word, mappings
}
