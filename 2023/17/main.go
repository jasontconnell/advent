package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func (pt xy) add(p2 xy) xy {
	return xy{pt.x + p2.x, pt.y + p2.y}
}

func (pt xy) mult(m int) xy {
	return xy{pt.x * m, pt.y * m}
}

type block struct {
	heatloss int
}

type state struct {
	pos      xy
	heatloss int
	dir      xy
	conslen  int
}

func (s state) String() string {
	return fmt.Sprintf("pos (%d, %d) dir (%d, %d): heatloss: %d conslen: %d", s.pos.x, s.pos.y, s.dir.x, s.dir.y, s.heatloss, s.conslen)
}

type key struct {
	pos     xy
	dir     xy
	conslen int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parseInput(in)
	mx, my := maxes(m)
	start := xy{0, 0}
	end := xy{mx, my}
	return traverse(m, start, end, 0, 3)
}

func part2(in input) output {
	m := parseInput(in)
	mx, my := maxes(m)
	start := xy{0, 0}
	end := xy{mx, my}
	return traverse(m, start, end, 4, 10)
}

func maxes(m map[xy]block) (int, int) {
	mx, my := 0, 0
	for k := range m {
		if k.x > mx {
			mx = k.x
		}
		if k.y > my {
			my = k.y
		}
	}
	return mx, my
}

func traverse(m map[xy]block, start xy, goal xy, minstraight, maxstraight int) int {
	queue := common.NewPriorityQueue(func(s state) float64 {
		return 1 / float64(s.heatloss)
	})

	itr := 0
	v := make(map[key]bool)
	min := math.MaxInt32
	queue.Enqueue(state{pos: start, heatloss: 0, dir: xy{1, 0}})
	queue.Enqueue(state{pos: start, heatloss: 0, dir: xy{0, 1}})
	for queue.Any() {
		cur := queue.Dequeue()

		k := key{pos: cur.pos, dir: cur.dir, conslen: cur.conslen}
		if _, ok := v[k]; ok {
			continue
		}
		v[k] = true

		if cur.pos == goal {
			if cur.heatloss < min {
				min = cur.heatloss
			}
			continue
		}

		mvs := getMoves(m, cur, goal, minstraight, maxstraight)
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
		itr++
	}
	return min
}

func getMoves(m map[xy]block, st state, goal xy, minstraight, maxstraight int) []state {
	dirs := []xy{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	potential := []xy{}

	if st.conslen == maxstraight {
		if st.dir.x != 0 {
			potential = append(potential, dirs[1], dirs[3])
		} else {
			potential = append(potential, dirs[0], dirs[2])
		}
	} else if minstraight > 0 && st.conslen < minstraight {
		potential = append(potential, st.dir)
	} else {
		for _, d := range dirs {
			if st.dir.x == -1 && d.x == 1 || st.dir.x == 1 && d.x == -1 ||
				st.dir.y == -1 && d.y == 1 || st.dir.y == 1 && d.y == -1 {
				continue
			}
			potential = append(potential, d)
		}
	}

	mvs := []state{}
	for _, d := range potential {
		dest := st.pos.add(d)
		b, ok := m[dest]
		if !ok {
			continue
		}

		if minstraight > 0 && d != st.dir { // is turning
			mt := d.mult(minstraight)
			minpt := st.pos.add(mt)

			if _, ok := m[minpt]; !ok {
				continue
			}
		}

		nst := state{pos: dest, heatloss: st.heatloss + b.heatloss, dir: d, conslen: 1}
		if st.dir == d {
			nst.conslen = st.conslen + 1
		}

		mvs = append(mvs, nst)
	}
	return mvs
}

func parseInput(in input) map[xy]block {
	m := make(map[xy]block)
	for y, line := range in {
		for x, c := range line {
			v, _ := strconv.Atoi(string(c))
			b := block{heatloss: v}
			m[xy{x, y}] = b
		}
	}
	return m
}
