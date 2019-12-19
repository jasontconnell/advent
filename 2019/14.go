package main

import (
	"bufio"
	"fmt"
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
	count int
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

			c, err := strconv.Atoi(r4[0])
			if err != nil {
				fmt.Println("error parsing", err)
			}

			name := r4[1]
			inchem := chemical{name: name, count: c}
			react.input = append(react.input, inchem)
		}

		r5 := strings.Fields(rs[1])
		c, err := strconv.Atoi(r5[0])
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

	fmt.Println("Time", time.Since(startTime))
}

func part1(reactions []reaction, fuelNeeded int) int {
	rmap := make(map[string]reaction)
	for _, r := range reactions {
		rmap[r.output.name] = r
	}

	needed := make(map[string]int)

	fuel, _ := rmap["FUEL"]
	provided := make(map[string]int)

	fulfill(fuel.output, 1, rmap, provided, needed)

	return provided["ORE"]
}

func fulfill(chem chemical, count int, rmap map[string]reaction, provided map[string]int, needed map[string]int) {
	r := rmap[chem.name]

	for _, i := range r.input {
		needed[i.name] += i.count
		re, ok := rmap[i.name]
		for provided[i.name] < needed[i.name] {
			fulfill(i, i.count, rmap, provided, needed)
			if ok {
				provided[i.name] += re.output.count
			} else {
				provided[i.name] += i.count
			}
		}
	}
}
