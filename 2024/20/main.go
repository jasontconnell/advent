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

func (p xy) dist(p2 xy) float64 {
	dx := p.x - p2.x
	dy := p.y - p2.y
	return math.Abs(float64(dx)) + math.Abs(float64(dy))
}

type state struct {
	pt         xy
	cheated    bool
	cheatstart xy
	cheatend   xy
	count      int
	cheatmove  bool
	dir        xy
	path       []xy
}

type statekey struct {
	pt         xy
	cheatstart xy
	cheatend   xy
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
	totalTime, optimalPath := traverse(m, start, end, false, 0, 0, nil, ex)
	countCheats, _ := traverse(m, start, end, true, totalTime, 100, optimalPath, ex)
	return countCheats
}

func part2(in input) output {
	return 0
}

func traverse(m map[xy]bool, start, end xy, allowCheat bool, compareValue, less int, optimalPath []xy, example bool) (int, []xy) {
	initial := state{pt: start, count: 0, cheated: false, cheatstart: xy{-1, -1}, cheatend: xy{-1, -1}}
	queue := common.NewPriorityQueue[state, float64](func(st state) float64 {
		return st.pt.dist(end)
	})
	queue.Enqueue(initial)
	visited := make(map[statekey]bool)
	total := 0
	minTime := math.MaxInt32
	var minpath []xy

	pmap := make(map[xy]int)
	if allowCheat && optimalPath != nil {
		for i, pt := range optimalPath {
			pmap[pt] = i
		}
	}

	leftmap := make(map[int]int)

	for queue.Any() {
		cur := queue.Dequeue()

		sk := statekey{pt: cur.pt, cheatstart: cur.cheatstart, cheatend: cur.cheatend}
		if _, ok := visited[sk]; ok {
			continue
		}
		visited[sk] = true

		if allowCheat && (compareValue-cur.count < less && !example) {
			continue
		}

		if cur.cheatmove {
			// skip to end
			eidx, eok := pmap[end]
			nstart, nok := pmap[cur.cheatend]
			if eok && nok {
				fullPath := cur.count + (eidx - nstart)
				saved := compareValue - fullPath
				if saved < 0 {
					continue
				}
				if saved >= less || example {
					total++
					leftmap[saved]++
				}
				continue
			}
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

		mvs := getMoves(m, cur, allowCheat)
		for _, mv := range mvs {
			if _, ok := m[mv.pt]; ok && mv.cheatmove {
				queue.Enqueue(mv)
			} else if ok {
				queue.Enqueue(mv)
			}
		}
	}
	if allowCheat {
		return total, nil
	}
	return minTime, minpath
}

func getMoves(m map[xy]bool, cur state, allowCheat bool) []state {
	mvs := []state{}
	dirs := []xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for _, d := range dirs {
		np := cur.pt.add(d)

		if open, ok := m[np]; ok {
			if allowCheat && !open && !cur.cheated {
				// np is a wall, check if we can go through a wall and end up on track
				cp := np.add(d)
				if f, ok := m[cp]; ok && f {
					mvs = append(mvs, state{pt: cp, dir: d, cheated: true, cheatstart: cur.pt, cheatend: cp, cheatmove: true, count: cur.count + 2})
				}
			}

			if open {
				var path []xy
				if !allowCheat {
					path = make([]xy, len(cur.path))
					copy(path, cur.path)
					path = append(path, np)
				}
				mvs = append(mvs, state{pt: np, path: path, dir: d, cheatstart: cur.cheatstart, cheatend: cur.cheatend, cheated: cur.cheated, cheatmove: false, count: cur.count + 1})
			}
		}
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
