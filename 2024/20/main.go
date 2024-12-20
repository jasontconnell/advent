package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

func (p xy) dist(p2 xy) int {
	dx := p.x - p2.x
	dy := p.y - p2.y
	return int(math.Abs(float64(dx)) + math.Abs(float64(dy)))
}

type state struct {
	pt                   xy
	count                int
	path                 []xy
	cheated              bool
	cheatdist            int
	cheatstart, cheatend xy
}

type statekey struct {
	pt                   xy
	cheatstart, cheatend xy
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m, start, end := parse(in)
	ex := len(in) < 20
	totalTime, optimalPath := traverse(m, start, end, 1, 0, 0, nil, ex)
	countCheats, _ := traverse(m, start, end, 2, totalTime, 100, optimalPath, ex)
	return countCheats
}

func part2(in input) output {
	m, start, end := parse(in)
	ex := len(in) < 20
	totalTime, optimalPath := traverse(m, start, end, 1, 0, 0, nil, ex)
	countCheats, _ := traverse(m, start, end, 20, totalTime, 100, optimalPath, ex)
	return countCheats
}

func traverse(m map[xy]bool, start, end xy, maxmoves int, compareValue, less int, optimalPath []xy, example bool) (int, []xy) {
	var minpath []xy

	minTime := math.MaxInt32
	lesstotal := 0
	allowCheat := maxmoves > 1

	pathMap := make(map[xy]int)
	if allowCheat && optimalPath != nil {
		for i, pt := range optimalPath {
			pathMap[pt] = i
		}
	}

	nullpoint := xy{-1, -1}

	initial := state{pt: start, count: 0, cheated: false, cheatdist: 0, cheatstart: nullpoint, cheatend: nullpoint, path: []xy{start}}
	queue := common.NewQueue[state, int]()
	visited := make(map[statekey]bool)

	queue.Enqueue(initial)
	for queue.Any() {
		cur := queue.Dequeue()

		if open, ok := m[cur.pt]; !ok || !open {
			continue
		}

		sk := statekey{pt: cur.pt, cheatstart: cur.cheatstart, cheatend: cur.cheatend}
		if _, ok := visited[sk]; ok {
			continue
		}
		visited[sk] = true

		if allowCheat && cur.cheatstart != nullpoint && cur.cheatend != nullpoint {
			// skip to end
			step := pathMap[cur.cheatend]
			fullPath := cur.count + (compareValue - step)
			saved := compareValue - fullPath
			if saved <= 0 && example {
				continue
			}
			if saved >= less || example {
				lesstotal++
			}
			continue
		}

		if cur.pt == end {
			if !allowCheat {
				if cur.count < minTime {
					minTime = cur.count
					minpath = cur.path
				}
			}
			continue
		}

		mvs := getMoves(pathMap, cur, maxmoves)
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
	}
	if allowCheat {
		return lesstotal, nil
	}
	return minTime, minpath
}

func getMoves(pathMap map[xy]int, cur state, maxsteps int) []state {
	mvs := []state{}

	if maxsteps == 1 {
		for _, d := range []xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			np := cur.pt.add(d)
			path := make([]xy, len(cur.path))
			copy(path, cur.path)
			path = append(path, np)
			mvs = append(mvs, state{pt: np, path: path, count: cur.count + 1})
		}
		return mvs
	}

	cidx := pathMap[cur.pt]
	// find moves that are > cidx and < maxsteps away
	for k, v := range pathMap {
		dist := k.dist(cur.pt)
		if dist > maxsteps || dist == 0 {
			continue
		}
		if v <= cidx {
			continue
		}
		if cur.cheated && dist > 1 {
			continue
		}

		cd := cur.cheatdist
		cs, ce := cur.cheatstart, cur.cheatend
		ch := cur.cheated
		if dist > 1 {
			cd = dist
			cs = cur.pt
			ce = k
			ch = true
		}

		mvs = append(mvs, state{pt: k, count: cur.count + k.dist(cur.pt), cheated: ch, cheatdist: cd, cheatstart: cs, cheatend: ce})
	}

	return mvs
}

func parse(in []string) (map[xy]bool, xy, xy) {
	m := make(map[xy]bool)
	var start, end xy

	for y := 0; y < len(in); y++ {
		for x := 0; x < len(in[y]); x++ {
			pt := xy{x, y}
			c := in[y][x]
			m[pt] = c != '#'

			if c == 'S' {
				start = pt
			} else if c == 'E' {
				end = pt
			}
		}
	}

	return m, start, end
}
