package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

var end xy = xy{31, 39}

type xy struct {
	x, y int
}

type state struct {
	pt    xy
	moves int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 13 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	x := parseInput(in)
	solve, moves, _ := navigate(xy{1, 1}, end, x, false, 0)
	if solve {
		return moves
	}
	return 0
}

func part2(in input) output {
	x := parseInput(in)
	_, _, unique := navigate(xy{1, 1}, end, x, true, 50)
	return unique
}

func parseInput(in input) int {
	i, _ := strconv.Atoi(in)
	return i
}

func navigate(start, goal xy, fav int, countUnique bool, uniqueWithin int) (bool, int, int) {
	visited := map[xy]bool{}
	queue := []state{}
	initial := state{pt: start, moves: 0}
	queue = append(queue, initial)
	solve := false
	moves := math.MaxInt32
	unique := 1

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.moves > moves {
			continue
		}

		if cur.pt == goal && !countUnique {
			solve = true
			if cur.moves < moves {
				moves = cur.moves
			}
			continue
		}

		if _, ok := visited[cur.pt]; ok {
			continue
		}
		visited[cur.pt] = true

		if cur.moves < uniqueWithin {
			unique++
		}

		mvs := getMoves(cur.pt, fav)
		for _, mv := range mvs {
			queue = append(queue, state{pt: mv, moves: cur.moves + 1})
		}
	}

	return solve, moves, unique
}

func getMoves(mv xy, fav int) (moves []xy) {
	for _, m := range []xy{{mv.x + 1, mv.y}, {mv.x, mv.y + 1}, {mv.x, mv.y - 1}, {mv.x - 1, mv.y}} {
		if m.x > -1 && m.y > -1 && !isWall(m, fav) {
			moves = append(moves, m)
		}
	}
	return
}

func isWall(mv xy, fav int) bool {
	x, y := mv.x, mv.y
	t := x*x + 3*x + 2*x*y + y + y*y
	t += fav
	binary := strconv.FormatInt(int64(t), 2)
	return len(strings.Replace(binary, "0", "", -1))%2 != 0
}
