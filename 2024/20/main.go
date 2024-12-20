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
	cheatstart xy
	cheatend   xy
	cheattime  int
	cheating   bool

	count int
	dir   xy
	path  []xy
}

type statekey struct {
	pt         xy
	cheatstart xy
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
	totalTime, optimalPath := traverse(m, start, end, 0, 0, 0, nil, ex)
	countCheats, _ := traverse(m, start, end, 2, totalTime, 100, optimalPath, ex)
	return countCheats
}

func part2(in input) output {
	return 0
	m, start, end := parse(in)
	ex := len(in) < 20
	totalTime, optimalPath := traverse(m, start, end, 0, 0, 0, nil, ex)
	countCheats, _ := traverse(m, start, end, 20, totalTime, 100, optimalPath, ex)
	return countCheats
}

func traverse(m map[xy]bool, start, end xy, cheattime int, compareValue, less int, optimalPath []xy, example bool) (int, []xy) {
	nullpoint := xy{-1, -1}
	initial := state{pt: start, count: 0, cheatstart: nullpoint, cheatend: nullpoint}
	queue := common.NewQueue[state, int]()
	queue.Enqueue(initial)
	visited := make(map[statekey]bool)
	total := 0
	minTime := math.MaxInt32
	var minpath []xy
	allowCheat := cheattime > 0

	cheatmap := make(map[xy]bool)
	for k := range m {
		cheatmap[k] = true
	}

	leftmap := make(map[int]int)

	pmap := make(map[xy]int)
	if allowCheat && optimalPath != nil {
		for i, pt := range optimalPath {
			pmap[pt] = i
		}
	}

	for queue.Any() {
		cur := queue.Dequeue()

		sk := statekey{pt: cur.pt, cheatstart: cur.cheatstart}
		if _, ok := visited[sk]; ok {
			continue
		}
		visited[sk] = true

		if allowCheat && (compareValue-cur.count < less && !example) {
			continue
		}

		if cur.cheatend != nullpoint {
			if open, ok := m[cur.cheatend]; !ok || !open {
				continue
			}
			// skip to end
			eidx, eok := pmap[end]
			nstart, nok := pmap[cur.cheatend]
			if eok && nok {
				fullPath := cur.count + (eidx - nstart)
				saved := compareValue - fullPath
				if saved <= 0 {
					continue
				}
				if saved >= less || example {
					leftmap[saved]++
					total++
				}
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

		var mvs []state
		if cur.cheating {
			mvs = getMoves(cheatmap, cur, allowCheat, cheattime)
		} else {
			mvs = getMoves(m, cur, allowCheat, cheattime)
		}
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
	}
	log.Println(leftmap)
	if allowCheat {
		return total, nil
	}
	return minTime, minpath
}

func getMoves(m map[xy]bool, cur state, allowCheat bool, cheattime int) []state {
	mvs := []state{}
	dirs := []xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for _, d := range dirs {
		np := cur.pt.add(d)

		if open, ok := m[np]; ok {
			if allowCheat && !open && !cur.cheating {
				mvs = append(mvs, state{pt: np, dir: d, cheattime: 1, cheatstart: np, cheatend: cur.cheatend, cheating: true, count: cur.count + 1})
			} else if allowCheat && open && cur.cheating && cur.cheattime >= cheattime {
				mvs = append(mvs, state{pt: np, dir: d, cheattime: cur.cheattime + 1, cheatstart: cur.cheatstart, cheatend: np, cheating: true, count: cur.count + 1})
			}

			if open {
				var path []xy
				if !allowCheat {
					path = make([]xy, len(cur.path))
					copy(path, cur.path)
					path = append(path, np)
				}

				if cur.cheating {
					mvs = append(mvs, state{pt: np, path: path, dir: d, cheatstart: cur.cheatstart, cheatend: np, cheattime: cur.cheattime + 1, count: cur.count + 1})
				} else {
					mvs = append(mvs, state{pt: np, path: path, dir: d, cheatstart: cur.cheatstart, cheatend: cur.cheatend, cheattime: cur.cheattime, count: cur.count + 1})
				}
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
