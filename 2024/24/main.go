package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int64

type wire struct {
	name  string
	id    int
	class rune
	value *int
}

type gate struct {
	str         string
	op          string
	left, right *wire
	output      *wire
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 24 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	w, g := parse(in)
	simulate(g)
	return calc(w, 'z')
}

func part2(in input) output {
	return 0
}

func simulate(gates []gate) {
	calced := true
	for calced {
		calced = false
		for _, g := range gates {
			if g.output.value == nil && g.left.value != nil && g.right.value != nil {
				v := doOp(g.left, g.right, g.op)
				g.output.value = &v
				calced = true
			}
		}
	}
}

func calc(wires map[string]*wire, class rune) int64 {
	var x int64
	for _, w := range wires {
		if w.class == 'z' && w.id >= 0 {
			if w.value == nil {
				panic("nil value")
			}
			if *w.value == 1 {
				log.Println(w.id, 1)
				x = x | (1 << w.id)
			}
		}
	}
	return x
}

func doOp(w1, w2 *wire, op string) int {
	v := 0
	switch op {
	case "AND":
		if *w1.value == *w2.value && *w2.value == 1 {
			v = 1
		}
	case "OR":
		if *w1.value == 1 || *w2.value == 1 {
			v = 1
		}
	case "XOR":
		if *w1.value != *w2.value {
			v = 1
		}
	}
	return v
}

func parse(in []string) (map[string]*wire, []gate) {
	wm := make(map[string]*wire)
	gates := []gate{}
	gstart := 0
	for i, line := range in {
		if line == "" {
			break
		}
		sp := strings.Fields(line)
		nm := sp[0]
		nm = nm[:len(nm)-1]
		val, _ := strconv.Atoi(sp[1])

		w := getOrCreate(nm, wm)
		w.value = &val

		gstart = i
	}
	for _, line := range in[gstart+2:] {
		sp := strings.Fields(line)
		w1name := sp[0]
		op := sp[1]
		w2name := sp[2]
		outname := sp[4]

		w1 := getOrCreate(w1name, wm)
		w2 := getOrCreate(w2name, wm)
		out := getOrCreate(outname, wm)

		g := gate{str: line, op: op, left: w1, right: w2, output: out}
		gates = append(gates, g)
	}
	return wm, gates
}

func getOrCreate(name string, wm map[string]*wire) *wire {
	var w *wire
	if t, ok := wm[name]; ok {
		w = t
	} else {
		w = &wire{name: name}
		id, err := strconv.Atoi(name[1:])
		if err != nil {
			id = -100
		}
		w.id = id
		w.class = rune(name[0])
		wm[name] = w
	}
	return w
}
