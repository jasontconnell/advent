package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "14.txt"

type reaction struct {
	output chemical
	input  []chemical
}

type chemical struct {
	name  string
	count int64
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	reactions := []reaction{}
	for scanner.Scan() {
		var txt = scanner.Text()
		rs := strings.Split(txt, "=>")
		r2 := strings.Split(rs[0], ",")

		react := reaction{}
		for _, r3 := range r2 {
			r4 := strings.Fields(r3)

			c, err := strconv.ParseInt(r4[0], 10, 64)
			if err != nil {
				fmt.Println("error parsing", err)
			}

			name := r4[1]
			inchem := chemical{name: name, count: c}
			react.input = append(react.input, inchem)
		}

		r5 := strings.Fields(rs[1])
		c, err := strconv.ParseInt(r5[0], 10, 64)
		if err != nil {
			fmt.Println("parsing out", err)
		}

		n := r5[1]

		react.output = chemical{name: n, count: c}

		reactions = append(reactions, react)
	}

	cp := make([]reaction, len(reactions))
	copy(cp, reactions)

	p1 := part1(cp, 1)
	fmt.Println("Part 1: ", p1)

	p2 := part2(cp, 1_000_000_000_000)
	fmt.Println("Part 2: ", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(reactions []reaction, fuelNeeded int64) int64 {
	rmap := make(map[string]reaction)
	for _, r := range reactions {
		rmap[r.output.name] = r
	}
	needed := make(map[string]int64)
	provided := make(map[string]int64)

	fuel := rmap["FUEL"]
	needed[fuel.output.name] = fuelNeeded

	fulfill(fuel.output, fuelNeeded, rmap, provided, needed)

	return provided["ORE"]
}

func part2(reactions []reaction, oreProvided int64) int64 {
	rmap := make(map[string]reaction)
	for _, r := range reactions {
		rmap[r.output.name] = r
	}

	done := false
	var low, high int64
	low = 0
	high = 20_000_000

	var result int64 = 0

	for !done {
		avail := make(map[string]int64)
		for _, r := range reactions {
			avail[r.output.name] = 0
		}

		avail["ORE"] = oreProvided
		start := (low + high) / 2

		made := makeN("FUEL", avail, rmap, start)

		if made > result {
			result = made
			low = start + 1
		} else {
			high = start - 1
		}

		m := makeN("FUEL", avail, rmap, 1)
		if m > 0 {
			result += m
		}
		done = low == high || (high-low) == 1
	}

	return result
}

func makeN(name string, avail map[string]int64, rmap map[string]reaction, n int64) int64 {
	reaction, ok := rmap[name]

	if !ok || len(reaction.input) == 0 {
		return avail[name]
	}

	madeAll := true
	for _, ichem := range reaction.input {
		chemr := rmap[ichem.name]
		output := chemr.output.count
		req := ichem.count * n

		if av, ok := avail[ichem.name]; ok && av < req {
			if output == 0 { // can't make any m'ore
				madeAll = false
				continue
			}
			numreq := int64(math.Ceil(float64(req) / float64(output)))
			made := makeN(ichem.name, avail, rmap, numreq)
			if made >= req {
				avail[ichem.name] += (made - req) // will be leftover, update for all immediately
			} else {
				madeAll = false
			}
		} else {
			avail[ichem.name] -= req
		}
	}

	if madeAll {
		return reaction.output.count * n
	}

	return 0
}

func fulfill(chem chemical, count int64, rmap map[string]reaction, provided map[string]int64, needed map[string]int64) {
	r := rmap[chem.name]
	for _, i := range r.input {
		needed[i.name] += i.count * count
		re, ok := rmap[i.name]
		for provided[i.name] < needed[i.name] {
			fulfill(i, count, rmap, provided, needed)
			if ok {
				provided[i.name] += re.output.count * count
			} else {
				provided[i.name] += i.count * count
			}
		}
	}
}
