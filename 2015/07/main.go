package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"
var reg *regexp.Regexp = regexp.MustCompile("^([a-z0-9]*?) ?([A-Z]*?) ?([a-z0-9]*?) -> ([a-z]{1,2})$")
var numreg *regexp.Regexp = regexp.MustCompile("^[0-9]+$")

type input []string

type Wire struct {
	Name string

	LeftEvaluated bool
	LeftString    string
	LeftWire      *Wire
	Left          uint16

	RightEvaluated bool
	RightString    string
	RightWire      *Wire
	Right          uint16

	Gate   string
	Result uint16

	Original string
}

func (w *Wire) String() string {
	return w.Name + ": { " + w.LeftString + " " + w.Gate + " " + w.RightString + " } ---- " + w.Original
}

func main() {
	startTime := time.Now()

	lines, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(lines)
	p2 := part2(lines, p1)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) uint16 {
	wm := getWires(in)

	return evaluate(wm, wm["a"])
}

func part2(in input, a uint16) uint16 {
	wm := getWires(in)
	wm["b"].Left = a
	return evaluate(wm, wm["a"])
}

func evalAll(wiremap map[string]*Wire) {
	for _, w := range wiremap {
		w.Result = evaluate(wiremap, w)
	}
}

func evaluate(wiremap map[string]*Wire, wire *Wire) uint16 {
	evalResult := uint16(0)
	if _, exists := wiremap[wire.Name]; exists {
		if wire.LeftWire != nil && !wire.LeftEvaluated {
			wire.Left = evaluate(wiremap, wire.LeftWire)
			wire.LeftEvaluated = true
		}

		if wire.RightWire != nil && !wire.RightEvaluated {
			wire.Right = evaluate(wiremap, wire.RightWire)
			wire.RightEvaluated = true
		}

		if wire.Gate != "" {
			evalResult = bitwise(wire.Left, wire.Right, wire.Gate)
		} else if wire.Left != 0 {
			evalResult = wire.Left
		} else if wire.Right != 0 {
			evalResult = wire.Right
		}

		wire.Result = evalResult
	}

	return evalResult
}

func bitwise(val1, val2 uint16, op string) uint16 {
	var ret uint16
	switch op {
	case "AND":
		ret = val1 & val2
		break
	case "OR":
		ret = val1 | val2
		break
	case "NOT":
		ret = ^val2
		break
	case "RSHIFT":
		ret = val1 >> val2
		break
	case "LSHIFT":
		ret = val1 << val2
	}

	return ret
}

func getWires(in input) map[string]*Wire {
	wiremap := make(map[string]*Wire)
	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			wire := &Wire{Name: groups[4], Gate: groups[2], Original: groups[0], LeftString: groups[1], RightString: groups[3]}

			if numreg.Match([]byte(wire.LeftString)) {
				if v1, perr := strconv.ParseUint(wire.LeftString, 10, 16); perr == nil {
					wire.Left = uint16(v1)
					wire.LeftEvaluated = true
				}
			}

			if numreg.Match([]byte(wire.RightString)) {
				if v2, perr := strconv.ParseUint(wire.RightString, 10, 16); perr == nil {
					wire.Right = uint16(v2)
					wire.RightEvaluated = true
				}
			}

			wiremap[wire.Name] = wire
		}
	}

	for _, w := range wiremap {
		if w.LeftString != "" {
			w.LeftWire = wiremap[w.LeftString]
		}

		if w.RightString != "" {
			w.RightWire = wiremap[w.RightString]
		}
	}

	return wiremap
}
