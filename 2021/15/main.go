package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type pqueue []state

func (pq pqueue) Len() int { return len(pq) }

func (pq pqueue) Less(i, j int) bool {
	return pq[i].score < pq[j].score
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	st := old[n-1]
	*pq = old[0 : n-1]
	return st
}

func (pq *pqueue) Push(x interface{}) {
	st := x.(state)
	*pq = append(*pq, st)
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

type grid [][]int
type state struct {
	pt    xy
	moves int
	level int
	score int
	dist  int
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
	fmt.Fprintln(w, "--2021 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	g := parseInput(in)
	start, goal := xy{0, 0}, xy{len(g[0]) - 1, len(g) - 1}
	solve, low := traverse(g, start, 0, goal, map[xy]int{})
	if !solve {
		return 0
	}
	return low
}

func part2(in input) output {
	g := parseInput(in)
	q := grow(g, 5)
	start, goal := xy{0, 0}, xy{len(q[0]) - 1, len(q) - 1}
	solve, low := traverse(q, start, 0, goal, map[xy]int{})
	if !solve {
		return 0
	}
	return low
}

func grow(g grid, q int) grid {
	quads := make(grid, len(g)*q)
	for i := 0; i < len(quads); i++ {
		quads[i] = make([]int, len(g[0])*q)
	}
	for ny := 0; ny < q; ny++ {
		for nx := 0; nx < q; nx++ {
			for y := 0; y < len(g); y++ {
				for x := 0; x < len(g[0]); x++ {
					yy := (ny * len(g)) + y
					xx := (nx * len(g[y])) + x
					quads[yy][xx] = g[y][x] + ny + nx

					if quads[yy][xx] >= 10 {
						quads[yy][xx] = (quads[yy][xx] % 10) + 1
					}
				}
			}
		}
	}
	return quads
}

func newState(g grid, pt xy, moves, level int) state {
	st := state{pt: pt}
	st.level = level
	st.moves = moves + 1
	st.dist = (len(g) - pt.y) + (len(g[0]) - pt.x)
	st.score = st.moves * st.dist
	return st
}

func traverse(g grid, start xy, startlevel int, goal xy, lows map[xy]int) (bool, int) {
	queue := pqueue{}
	initial := newState(g, start, 0, startlevel)
	queue.Push(initial)
	heap.Init(&queue)
	lowest := math.MaxInt32
	solve := false
	for queue.Len() > 0 {
		cur := (heap.Pop(&queue)).(state)
		if cur.pt == goal {
			solve = true
			if cur.level < lowest {
				lowest = cur.level
			}
			continue
		}

		if cur.level > lowest {
			continue
		}

		mvs := getMoves(g, cur.pt)
		for _, mv := range mvs {
			nlv := cur.level + g[mv.y][mv.x]
			n := newState(g, mv, cur.moves, nlv)
			if lv, ok := lows[mv]; ok && nlv >= lv {
				continue
			}
			lows[mv] = nlv
			queue = append(queue, n)
		}
	}
	return solve, lowest
}

func getMoves(g grid, cur xy) []xy {
	pts := []xy{{cur.x + 1, cur.y}, {cur.x - 1, cur.y}, {cur.x, cur.y + 1}, {cur.x, cur.y - 1}}

	valid := []xy{}
	for _, pt := range pts {
		if pt.x < 0 || pt.x == len(g[0]) || pt.y < 0 || pt.y == len(g) {
			continue
		}
		valid = append(valid, pt)
	}

	return valid
}

func parseInput(in input) grid {
	ret := make(grid, len(in))

	for i, line := range in {
		for _, c := range line {
			v, _ := strconv.Atoi(string(c))

			ret[i] = append(ret[i], v)
		}
	}
	return ret
}
