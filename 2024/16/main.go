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
	return traverse(m, start, end)
}

func part2(in input) output {
	return 0
}

func traverse(m map[xy]rune, start, end xy) int {
	queue := common.NewPriorityQueue(func(s state) int {
		return s.score
	})
	queue.Enqueue(state{pt: start, score: 0, facing: Right})

	visit := make(map[statekey]state)

	lowscore := math.MaxInt32

	for queue.Any() {
		cur := queue.Dequeue()

		if cur.score > lowscore {
			continue
		}

		sk := statekey{pt: cur.pt, facing: cur.facing}
		if s, ok := visit[sk]; ok {
			if cur.score >= s.score {
				continue
			}
		}
		visit[sk] = cur

		if cur.pt == end {
			if cur.score < lowscore {
				lowscore = cur.score
			}
			continue
		}

		mvs := getMoves(cur)
		for _, mv := range mvs {
			if _, ok := m[mv.pt]; ok {
				queue.Enqueue(mv)
			}
		}
	}
	return lowscore
}

func getMoves(cur state) []state {
	mvs := []state{}
	mvs = append(mvs, state{pt: cur.pt.add(cur.facing.xy()), score: cur.score + 1, facing: cur.facing})

	sd := cur.facing
	switch sd {
	case Up, Down:
		mvs = append(mvs, state{pt: cur.pt, score: cur.score + 1000, facing: Right})
		mvs = append(mvs, state{pt: cur.pt, score: cur.score + 1000, facing: Left})
		// rotate twice to get opposite
		if sd == Up {
			mvs = append(mvs, state{pt: cur.pt, score: cur.score + 2000, facing: Down})
		} else {
			mvs = append(mvs, state{pt: cur.pt, score: cur.score + 2000, facing: Up})
		}
	case Left, Right:
		mvs = append(mvs, state{pt: cur.pt, score: cur.score + 1000, facing: Up})
		mvs = append(mvs, state{pt: cur.pt, score: cur.score + 1000, facing: Down})
		if sd == Left {
			mvs = append(mvs, state{pt: cur.pt, score: cur.score + 2000, facing: Right})
		} else {
			mvs = append(mvs, state{pt: cur.pt, score: cur.score + 2000, facing: Left})
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
