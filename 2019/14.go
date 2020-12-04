package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "14_test.txt"

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

	p2 := part2(cp, p1) // 1_000_000_000_000)
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

	// needed := make(map[string]int64)
	provided := make(map[string]int64)

	provided["ORE"] = oreProvided
	provided["FUEL"] = 0

	var total int64
	pfuel := int64(1)
	for provided["ORE"] > 0 {
		total += pfuel
		fulfill2(rmap["FUEL"].output, pfuel, rmap, provided)
		if provided["ORE"]%100000 == 0 {
			fmt.Println(provided["ORE"])
		}

		if provided["ORE"] < 1000 {
			fmt.Println("less then 1 thousand", provided["ORE"])
			pfuel = 1
		} else if provided["ORE"] < 100_000 {
			fmt.Println("less then 100 thousand", provided["ORE"])
			pfuel = 1
		} else if provided["ORE"] < 100_000_000 {
			fmt.Println("less then 10 million", provided["ORE"])
			pfuel = 100
		}
	}

	return total
}

func fulfill2(chem chemical, count int64, rmap map[string]reaction, provided map[string]int64) {
	// fmt.Println(chem.name, count)
	r := rmap[chem.name]
	for _, i := range r.input {
		re, ok := rmap[i.name]
		var c int64 = count
		for provided[i.name] < count {
			fulfill2(i, i.count*count, rmap, provided)
			if ok {
				c = re.output.count * count
				provided[i.name] += c
			} else {
				c = i.count * count
				provided[i.name] += c
			}
		}
		provided[i.name] -= c
	}
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
