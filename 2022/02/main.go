package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type strat struct {
	move    rune
	counter rune
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
	fmt.Fprintln(w, "--2022 day 02 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	strats := parseInput(in)
	sum := 0
	for _, st := range strats {
		sum += getResult(st)
	}
	return sum
}

func part2(in input) output {
	return 0
}

func getResult(st strat) int {
	score := 0
	mvpts := 0
	switch st.counter {
	case 'X':
		mvpts = 1
		if st.move == 'A' {
			score = 3
		} else if st.move == 'C' {
			score = 6
		}
	case 'Y':
		mvpts = 2
		if st.move == 'A' {
			score = 6
		} else if st.move == 'B' {
			score = 3
		}
	case 'Z':
		mvpts = 3
		if st.move == 'B' {
			score = 6
		} else if st.move == 'C' {
			score = 3
		}
	}
	return score + mvpts
}

func parseInput(lines input) []strat {
	strats := []strat{}
	for _, line := range lines {
		sp := strings.Split(line, " ")
		st := strat{move: rune(sp[0][0]), counter: rune(sp[1][0])}
		strats = append(strats, st)
	}
	return strats
}
