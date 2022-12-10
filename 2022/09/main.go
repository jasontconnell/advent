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

func print(m map[xy]int) {
	minx, maxx := math.MaxInt32, math.MinInt32
	miny, maxy := math.MaxInt32, math.MinInt32
	for k := range m {
		if k.x > maxx {
			maxx = k.x
		}
		if k.x < minx {
			minx = k.x
		}
		if k.y > maxy {
			maxy = k.y
		}
		if k.y < miny {
			miny = k.y
		}
	}

	for y := miny; y < maxy+1; y++ {
		for x := minx; x < maxx+1; x++ {
			if _, ok := m[xy{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
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
	return traverse(moves, 2)
}

func part2(in input) output {
	moves := parseInput(in)
	return traverse(moves, 10)
}

func traverse(moves []move, n int) int {
	pos := make([]xy, n)
	visit := make([]map[xy]int, n)
	for i := 0; i < n; i++ {
		pos[i] = xy{}
		visit[i] = make(map[xy]int)
		visit[i][pos[i]] = 1
	}

	for _, mv := range moves {
		for i := 0; i < mv.num; i++ {
			moveOne(mv.dir, pos, visit)
		}
	}

	return len(visit[n-1])
}

func moveOne(d dir, pos []xy, visit []map[xy]int) {
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

	pos[0].x += dx
	pos[0].y += dy
	visit[0][pos[0]]++

	for i := 1; i < len(pos); i++ {
		fpos, bpos := pos[i-1], pos[i]
		if !touching(fpos, bpos) {
			tdx, tdy := 0, 0

			if bpos.x > fpos.x {
				tdx = -1
			} else if bpos.x < fpos.x {
				tdx = 1
			}

			if bpos.y < fpos.y {
				tdy = 1
			} else if bpos.y > fpos.y {
				tdy = -1
			}

			pos[i].x += tdx
			pos[i].y += tdy
			visit[i][pos[i]]++
		}
	}
}

func touching(p1, p2 xy) bool {
	y := math.Abs(float64(p2.y - p1.y))
	x := math.Abs(float64(p2.x - p1.x))
	if y > 1 || x > 1 { //1,1 is diagonal but still touching
		return false
	}
	return true
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
