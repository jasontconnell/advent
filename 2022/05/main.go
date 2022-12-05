package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = string

type stack struct {
	boxes []rune
}

func (s stack) String() string {
	str := ""
	for _, r := range s.boxes {
		str += "[" + string(r) + "] "
	}
	return str
}

type move struct {
	from, to int
	quantity int
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
	fmt.Fprintln(w, "--2022 day 05 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	stacks, moves := parseInput(in)
	return getTops(stacks, moves, false)
}

func part2(in input) output {
	stacks, moves := parseInput(in)
	return getTops(stacks, moves, true) // use CrateMover9001
}

func getTops(stacks []stack, moves []move, is9001 bool) string {
	for _, mv := range moves {
		if !is9001 {
			for i := 0; i < mv.quantity; i++ {
				stacks = moveOne(stacks, mv.from, mv.to)
			}
		} else {
			stacks = moveMultiple(stacks, mv.from, mv.to, mv.quantity)
		}
	}

	str := ""
	for _, s := range stacks {
		slen := len(s.boxes)
		str += string(s.boxes[slen-1])
	}
	return str
}

func moveOne(stacks []stack, from, to int) []stack {
	flen := len(stacks[from].boxes)

	r := stacks[from].boxes[flen-1]
	stacks[from].boxes = stacks[from].boxes[:flen-1]
	stacks[to].boxes = append(stacks[to].boxes, r)
	return stacks
}

func moveMultiple(stacks []stack, from, to, quantity int) []stack {
	flen := len(stacks[from].boxes)

	start := flen - quantity
	sstk := stacks[from].boxes[start:]
	stacks[from].boxes = stacks[from].boxes[:start]
	stacks[to].boxes = append(stacks[to].boxes, sstk...)
	return stacks
}

func parseInput(in input) ([]stack, []move) {
	num := len(in[0])/4 + 1

	stacks := make([]stack, num)

	movestart := 0
	for i, line := range in {
		if len(line) == 0 {
			movestart = i
			break
		}

		for c := 0; c < num; c++ {
			r := rune(line[c*4+1])
			if r != ' ' && unicode.IsLetter(r) {
				stacks[c].boxes = append([]rune{r}, stacks[c].boxes...)
			}
		}
	}

	moves := []move{}
	reg := regexp.MustCompile("move ([0-9]+) from ([0-9]+) to ([0-9]+)")

	for i := movestart; i < len(in); i++ {
		if groups := reg.FindStringSubmatch(in[i]); groups != nil && len(groups) > 1 {
			q, _ := strconv.Atoi(groups[1])
			f, _ := strconv.Atoi(groups[2])
			t, _ := strconv.Atoi(groups[3])
			moves = append(moves, move{quantity: q, from: f - 1, to: t - 1}) // convert to index
		}
	}

	return stacks, moves
}
