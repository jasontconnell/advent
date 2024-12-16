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

type dir int

const (
	Up dir = iota
	Down
	Left
	Right
)

func (d dir) xy() xy {
	var p xy
	switch d {
	case Up:
		p.y = -1
	case Down:
		p.y = 1
	case Left:
		p.x = -1
	case Right:
		p.x = 1
	}
	return p
}

type state struct {
	pt     xy
	score  int
	facing dir
	path   []state
	isturn bool
}

type statekey struct {
	pt     xy
	facing dir
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 16 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m, start, end := parse(in)
	lowscore, _ := traverse(m, start, end, false, math.MaxInt32)
	return lowscore
}

func part2(in input) output {
	m, start, end := parse(in)
	lowscore, _ := traverse(m, start, end, false, math.MaxInt32)
	_, length := traverse(m, start, end, true, lowscore)
	return length
}

func traverse(m map[xy]rune, start, end xy, trackpath bool, maxscore int) (int, int) {
	queue := common.NewPriorityQueue(func(s state) int {
		return s.score
	})
	startstate := state{pt: start, score: 0, facing: Right}
	startstate.path = append(startstate.path, startstate)
	queue.Enqueue(startstate)

	visit := make(map[statekey]state)
	best := make(map[xy]state)

	lowscore := math.MaxInt32
	solves := 0

	for queue.Any() {
		cur := queue.Dequeue()

		if cur.score > lowscore {
			continue
		}

		if trackpath && cur.score > maxscore+1 {
			continue
		}

		sk := statekey{pt: cur.pt, facing: cur.facing}
		if s, ok := visit[sk]; ok {
			if !trackpath && cur.score >= s.score {
				continue
			} else if trackpath {
				b, ok := best[cur.pt]
				if ok && cur.score > b.score || cur.score > s.score {
					continue
				}
				same := true
				if len(cur.path) == len(s.path) {
					for i := 0; i < len(cur.path); i++ {
						if cur.path[i].pt != s.path[i].pt {
							same = false
							break
						}
					}
				}
				if same {
					continue
				}
			}
		}
		visit[sk] = cur

		if cur.pt == end {
			if cur.score < lowscore && trackpath {
				for k := range best {
					delete(best, k)
				}
			}
			if cur.score <= lowscore {
				solves++
				lowscore = cur.score
				if trackpath {
					for _, st := range cur.path {
						best[st.pt] = st
					}
				}
			}

			continue
		}

		mvs := getMoves(cur, trackpath)
		for _, mv := range mvs {
			if _, ok := m[mv.pt]; ok {
				queue.Enqueue(mv)
			}
		}
	}
	return lowscore, len(best)
}

func getMoves(cur state, trackpath bool) []state {
	var pcopy []state

	next := cur.pt.add(cur.facing.xy())

	if trackpath {
		pcopy = make([]state, len(cur.path))
		copy(pcopy, cur.path)
	}

	mvs := []state{}
	mvstate := state{pt: next, score: cur.score + 1, facing: cur.facing}
	if trackpath {
		mvstate.path = append(pcopy, mvstate)
	}
	mvs = append(mvs, mvstate)

	if cur.isturn {
		return mvs
	}

	switch cur.facing {
	case Up, Down:
		state1 := state{pt: cur.pt, score: cur.score + 1000, facing: Right, isturn: true}
		state2 := state{pt: cur.pt, score: cur.score + 1000, facing: Left, isturn: true}
		state3 := state{pt: cur.pt, score: cur.score + 2000, facing: Down, isturn: true}
		state4 := state{pt: cur.pt, score: cur.score + 2000, facing: Up, isturn: true}
		if trackpath {
			state1.path = append(pcopy, state1)
			state2.path = append(pcopy, state2)
		}
		mvs = append(mvs, state1)
		mvs = append(mvs, state2)
		// rotate twice to get opposite
		if cur.facing == Up {
			state3.path = append(pcopy, state3)
			mvs = append(mvs, state3)
		} else {
			state4.path = append(pcopy, state4)
			mvs = append(mvs, state4)
		}
	case Left, Right:
		state1 := state{pt: cur.pt, score: cur.score + 1000, facing: Up, isturn: true}
		state2 := state{pt: cur.pt, score: cur.score + 1000, facing: Down, isturn: true}
		state3 := state{pt: cur.pt, score: cur.score + 2000, facing: Right, isturn: true}
		state4 := state{pt: cur.pt, score: cur.score + 2000, facing: Left, isturn: true}
		if trackpath {
			state1.path = append(pcopy, state1)
			state2.path = append(pcopy, state2)
		}
		mvs = append(mvs, state1)
		mvs = append(mvs, state2)
		// rotate twice to get opposite
		if cur.facing == Left {
			state3.path = append(pcopy, state3)
			mvs = append(mvs, state3)
		} else {
			mvs = append(mvs, state4)
			state4.path = append(pcopy, state4)
		}

	}
	return mvs
}

func parse(in []string) (map[xy]rune, xy, xy) {
	m := make(map[xy]rune)
	var start, end xy
	for y, line := range in {
		for x, c := range line {
			if c == '#' {
				continue
			}
			pt := xy{x, y}
			if c == 'S' {
				start = pt
			} else if c == 'E' {
				end = pt
			}
			m[pt] = c
		}
	}
	return m, start, end
}
