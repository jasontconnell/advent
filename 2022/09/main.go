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

type input = []string
type output = int

type dir string

const (
	U dir = "U"
	D dir = "D"
	L dir = "L"
	R dir = "R"
)

type move struct {
	dir dir
	num int
}

type xy struct {
	x, y int
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
	fmt.Fprintln(w, "--2022 day 09 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	moves := parseInput(in)
	head, tail := make(map[xy]int), make(map[xy]int)
	head[xy{0, 0}] = 1
	tail[xy{0, 0}] = 1
	traverse(moves, head, tail)
	return len(tail)
}

func part2(in input) output {
	return 0
}

func traverse(moves []move, head, tail map[xy]int) {
	hpos, tpos := xy{0, 0}, xy{0, 0}
	for _, mv := range moves {
		for i := 0; i < mv.num; i++ {
			hpos, tpos = moveOne(mv.dir, hpos, tpos, head, tail)
		}
	}
}

func moveOne(d dir, hpos, tpos xy, head, tail map[xy]int) (xy, xy) {
	var dx, dy int
	switch d {
	case U:
		dy = 1
	case D:
		dy = -1
	case L:
		dx = -1
	case R:
		dx = 1
	}

	hpos.x += dx
	hpos.y += dy
	head[hpos]++

	// not touching, move tail
	if dist(hpos, tpos) > 1 {
		tdx, tdy := 0, 0
		switch d {
		case U:
			tdy = 1
		case D:
			tdy = -1
		case L:
			tdx = -1
		case R:
			tdx = 1
		}

		// take care of diagonal
		switch d {
		case U, D:
			if tpos.x < hpos.x {
				tdx = 1
			} else if tpos.x > hpos.x {
				tdx = -1
			}
		case R, L:
			if tpos.y < hpos.y {
				tdy = 1
			} else if tpos.y > hpos.y {
				tdy = -1
			}
		}

		tpos.x += tdx
		tpos.y += tdy
		tail[tpos]++
	}

	return hpos, tpos
}

func dist(p1, p2 xy) int {
	y := math.Abs(float64(p2.y - p1.y))
	x := math.Abs(float64(p2.x - p1.x))
	if y > 1 || x > 1 { //1,1 is diagonal but still touching
		return int(y + x)
	}
	return 1
}

func parseInput(in input) []move {
	mvs := []move{}
	for _, line := range in {
		sp := strings.Split(line, " ")
		if len(sp) == 2 {
			d := dir(sp[0])
			n, _ := strconv.Atoi(sp[1])
			mvs = append(mvs, move{dir: d, num: n})
		}
	}
	return mvs
}
