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
type output = uint64

type space struct {
	val  int
	next *space
}

type player struct {
	id     int
	pos    *space
	score  int
	posval int
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
	_, loser, turns := play([]*player{p1, p2}, 1000, 100, 3)
	return uint64(loser.score * turns)
}

func part2(in input) output {
	board := getSpaces(10)
	p1, p2 := parseInput(in, board)
	return playQuantum(p1, p2, 10, 21, 3)
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

type qgame struct {
	curpos, otherpos     int
	curscore, otherscore int
}

type qres struct {
	cur, other uint64
}

func playQuantum(p1, p2 *player, boardsize, winscore, dmax int) uint64 {
	game := qgame{curpos: p1.posval - 1, otherpos: p2.posval - 1, curscore: 0, otherscore: 0}
	m := make(map[qgame]qres)
	res := quantumGame(m, game, boardsize, winscore, dmax)

	ret := res.other
	if res.cur > res.other {
		ret = res.cur
	}
	return ret
}

func quantumGame(m map[qgame]qres, game qgame, boardsize, winscore, dmax int) qres {
	if game.curscore >= 21 {
		return qres{1, 0}
	}
	if game.otherscore >= 21 {
		return qres{0, 1}
	}
	if v, ok := m[game]; ok {
		return v
	}
	var final qres
	for d1 := 1; d1 <= dmax; d1++ {
		for d2 := 1; d2 <= dmax; d2++ {
			for d3 := 1; d3 <= dmax; d3++ {
				rsum := d1 + d2 + d3

				npos := (game.curpos + rsum) % boardsize
				nscore := game.curscore + npos + 1
				g := qgame{curpos: game.otherpos, curscore: game.otherscore, otherpos: npos, otherscore: nscore}

				res := quantumGame(m, g, boardsize, winscore, dmax)
				final.cur += res.other
				final.other += res.cur
			}
		}
	}
	m[game] = final
	return final
}

func play(players []*player, winscore, dmax, rolls int) (*player, *player, int) {
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

		if p.score >= winscore {
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

	p1pos, _ := strconv.Atoi(p1sp[len(p1sp)-1])
	p2pos, _ := strconv.Atoi(p2sp[len(p1sp)-1])
	p1 := player{id: 1, posval: p1pos}
	p2 := player{id: 2, posval: p2pos}

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
