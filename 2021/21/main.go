package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type space struct {
	val  int
	next *space
}

type player struct {
	id    int
	pos   *space
	score int
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
	fmt.Fprintln(w, "--2021 day 21 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	board := getSpaces(10)
	p1, p2 := parseInput(in, board)
	_, loser, turns := play([]*player{p1, p2}, 100, 3)
	fmt.Println(loser.score, turns)
	return loser.score * turns
}

func part2(in input) output {
	return 0
}

func getSpaces(max int) *space {
	root := &space{val: 1}
	p := root
	var last *space
	for i := 2; i < max+1; i++ {
		n := &space{val: i}
		p.next = n
		if i == max {
			last = n
		}
		p = p.next
	}

	last.next = root

	return root
}

func play(players []*player, dmax int, rolls int) (*player, *player, int) {
	dcur := 1
	pidx := 0
	done := false
	var winner, loser *player
	turns := 0

	for !done {
		p := players[pidx]

		for r := 0; r < rolls; r++ {
			for i := 0; i < dcur; i++ {
				p.pos = p.pos.next
			}
			dcur = (dcur + 1) % (dmax + 1)
			if dcur == 0 {
				dcur = 1
			}
			turns++
		}

		p.score += p.pos.val

		if p.score >= 1000 {
			winner = p
			loser = players[(pidx+1)%2]
			done = true
		}

		pidx = (pidx + 1) % 2

	}

	return winner, loser, turns
}

func parseInput(in input, board *space) (*player, *player) {
	p1sp := strings.Fields(in[0])
	p2sp := strings.Fields(in[1])

	p1 := player{id: 1}
	p2 := player{id: 2}

	p1pos, _ := strconv.Atoi(p1sp[len(p1sp)-1])
	p2pos, _ := strconv.Atoi(p2sp[len(p1sp)-1])

	r := board

	for p1.pos == nil || p2.pos == nil {
		if r.val == p1pos {
			p1.pos = r
		}

		if r.val == p2pos {
			p2.pos = r
		}
		r = r.next
	}

	return &p1, &p2
}
