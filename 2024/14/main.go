package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type bot struct {
	pos xy
	v   xy
}
type xy struct {
	x, y int
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	bots := parse(in)
	w, h := 101, 103
	if len(bots) < 20 {
		w = 11
		h = 7
	}
	bots = simulate(bots, 100, w, h)
	return countQuadrants(bots, w, h)
}

func part2(in input) output {
	bots := parse(in)
	w, h := 101, 103
	if len(bots) < 20 {
		w = 11
		h = 7
	}
	return findTree(bots, w, h)
}

func findTree(bots []bot, w, h int) int {
	n := 1
	for {
		bots = simulate(bots, 1, w, h)
		if rootMeanSquare(bots) < 42 {
			print(bots, w, h)
			break
		}
		n++
	}
	return n
}

func print(bots []bot, w, h int) {
	m := make(map[xy]bot)
	for _, b := range bots {
		m[b.pos] = b
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			pt := xy{x, y}
			if _, ok := m[pt]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func rootMeanSquare(bots []bot) float64 {
	var ans int
	for i := 0; i < len(bots); i++ {
		for j := 0; j < len(bots); j++ {
			if i == j {
				continue
			}
			x1, y1 := bots[i].pos.x, bots[i].pos.y
			x2, y2 := bots[j].pos.x, bots[j].pos.y

			ans += (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
		}
	}
	fnum := float64(len(bots))
	return math.Sqrt(float64(ans) / (fnum * fnum))
}

func simulate(bots []bot, sec, w, h int) []bot {
	for i := 0; i < sec; i++ {
		for idx, b := range bots {
			nk := b.pos.add(b.v)
			if nk.x >= w {
				nk.x = nk.x % w
			} else if nk.x < 0 {
				nk.x = w + nk.x
			}
			if nk.y >= h {
				nk.y = nk.y % h
			} else if nk.y < 0 {
				nk.y = h + nk.y
			}
			bots[idx].pos = nk
		}
	}
	return bots
}

func countQuadrants(bots []bot, w, h int) int {
	quads := make(map[int]int)
	for _, b := range bots {
		k := b.pos
		if k.x < w/2 && k.y < h/2 {
			quads[0]++
		} else if k.x > w/2 && k.y < h/2 {
			quads[1]++
		} else if k.x < w/2 && k.y > h/2 {
			quads[2]++
		} else if k.x > w/2 && k.y > h/2 {
			quads[3]++
		}
	}
	total := 1
	for _, v := range quads {
		total *= v
	}
	return total
}

func parse(in []string) []bot {
	reg := regexp.MustCompile(`^p=([0-9]+),([0-9]+) v=(-?[0-9]+),(-?[0-9]+)$`)

	bots := []bot{}
	for _, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) != 5 {
			continue
		}
		px, _ := strconv.Atoi(m[1])
		py, _ := strconv.Atoi(m[2])

		vx, _ := strconv.Atoi(m[3])
		vy, _ := strconv.Atoi(m[4])

		b := bot{pos: xy{px, py}, v: xy{vx, vy}}
		bots = append(bots, b)
	}
	return bots
}
