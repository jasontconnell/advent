package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

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

type state struct {
	step int
	pt   xy
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	list := parse(in)
	steps := 1024
	s := 70
	if len(list) == 25 {
		s = 6
		steps = 12
	}
	c, _ := traverse(list, steps, s)
	return c
}

func part2(in input) string {
	list := parse(in)
	steps := 1024
	s := 70
	if len(list) == 25 {
		s = 6
		steps = 12
	}
	p := getBlockingPoint(steps, list, s)
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func getBlockingPoint(start int, bytes []xy, gridsize int) xy {
	var p xy
	min, max := start, len(bytes)
	if len(bytes)%2 != 0 {
		max--
	}
	for {
		mid := (min + max) / 2
		_, found := traverse(bytes, mid, gridsize)
		if found {
			min = mid + 1
		} else {
			max = mid - 1
		}
		if max == min {
			p = bytes[mid+1] // the next point is the breaking point
			break
		}
	}
	return p
}

func traverse(bytes []xy, numbytes, gridsize int) (int, bool) {
	goal := xy{gridsize, gridsize}
	queue := common.NewQueue[state, int]()
	queue.Enqueue(state{pt: xy{0, 0}, step: 0})
	visited := make(map[xy]bool)
	lowest := math.MaxInt32
	found := false

	for queue.Any() {
		cur := queue.Dequeue()

		if _, ok := visited[cur.pt]; ok {
			continue
		}
		visited[cur.pt] = true

		if cur.pt == goal {
			if cur.step < lowest {
				lowest = cur.step
			}
			found = true
			continue
		}

		mvs := getMoves(cur, bytes, gridsize, numbytes)
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
	}
	return lowest, found
}

func getMoves(st state, bytes []xy, gridsize int, maxbytes int) []state {
	mvs := []state{}
	dirs := []xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	g := getGridAt(bytes, maxbytes)
	for _, d := range dirs {
		p := st.pt.add(d)
		inbounds := p.x >= 0 && p.x <= gridsize && p.y >= 0 && p.y <= gridsize
		if _, ok := g[p]; !ok && inbounds {
			nstate := state{step: st.step + 1, pt: p}
			mvs = append(mvs, nstate)
		}
	}
	return mvs
}

func getGridAt(bytes []xy, maxbytes int) map[xy]bool {
	g := make(map[xy]bool)
	for i := 0; i < maxbytes; i++ {
		g[bytes[i]] = true
	}
	return g
}

func parse(in []string) []xy {
	list := []xy{}
	for _, line := range in {
		sp := strings.Split(line, ",")
		if len(sp) == 2 {
			x, _ := strconv.Atoi(sp[0])
			y, _ := strconv.Atoi(sp[1])

			list = append(list, xy{x, y})
		}
	}
	return list
}
