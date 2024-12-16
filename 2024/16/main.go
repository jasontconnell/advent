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
	path   []xyscore
	isturn bool
}

type statekey struct {
	pt     xy
	facing dir
}

type xyscore struct {
	pt    xy
	score int
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
	lowscore, _ := traverse(m, start, end, false)
	return lowscore
}

func part2(in input) output {
	m, start, end := parse(in)
	_, length := traverse(m, start, end, true)
	return length
}

func traverse(m map[xy]rune, start, end xy, trackpath bool) (int, int) {
	queue := common.NewPriorityQueue(func(s state) int {
		h := s.score
		if trackpath {
			h = h * len(s.path)
		}
		return h
	})
	startstate := state{pt: start, score: 0, facing: Right}
	startstate.path = append(startstate.path, xyscore{pt: start, score: 0})
	queue.Enqueue(startstate)

	visit := make(map[statekey]state)
	best := make(map[xy]xyscore)
	lowscore := math.MaxInt32

	for queue.Any() {
		cur := queue.Dequeue()

		if cur.score > lowscore {
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
			if trackpath && cur.score < lowscore {
				best = make(map[xy]xyscore)
			}
			if cur.score <= lowscore {
				lowscore = cur.score
				if trackpath {
					for _, st := range cur.path {
						best[st.pt] = st
					}
				}
			}

			continue
		}

		mvs := getMoves(m, cur, trackpath)
		for _, mv := range mvs {
			if _, ok := m[mv.pt]; ok {
				queue.Enqueue(mv)
			}
		}
	}
	return lowscore, len(best)
}

func getMoves(m map[xy]rune, cur state, trackpath bool) []state {
	var pcopy []xyscore

	if trackpath {
		pcopy = make([]xyscore, len(cur.path))
		copy(pcopy, cur.path)
	}

	mvs := []state{}

	next := cur.pt.add(cur.facing.xy())
	mvstate := state{pt: next, score: cur.score + 1, facing: cur.facing}
	if trackpath {
		mvstate.path = append(pcopy, xyscore{pt: mvstate.pt, score: mvstate.score})
	}
	mvs = append(mvs, mvstate)

	if cur.isturn {
		return mvs
	}

	right := state{pt: cur.pt, score: cur.score, facing: Right, isturn: true}
	left := state{pt: cur.pt, score: cur.score, facing: Left, isturn: true}
	down := state{pt: cur.pt, score: cur.score, facing: Down, isturn: true}
	up := state{pt: cur.pt, score: cur.score, facing: Up, isturn: true}

	tmvs := []state{right, left, down, up}
	for i := 0; i < len(tmvs); i++ {
		s := tmvs[i]
		fpt := s.pt.add(s.facing.xy())
		if _, ok := m[fpt]; !ok {
			continue
		}
		if s.facing == cur.facing {
			continue
		}
		add := false
		switch cur.facing {
		case Up, Down:
			if s.facing == Left || s.facing == Right {
				s.score += 1000
				add = true
			}
		case Left, Right:
			if s.facing == Up || s.facing == Down {
				s.score += 1000
				add = true
			}
		}
		s.score += 1
		s.pt = fpt
		if trackpath {
			s.path = append(pcopy, xyscore{pt: s.pt, score: s.score})
		}
		if add {
			mvs = append(mvs, s)
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
