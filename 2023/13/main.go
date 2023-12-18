package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type grid map[xy]terrain

func (g grid) ItemColRow(x, y int) terrain {
	pt := xy{x, y}
	return g[pt]
}

func (g grid) ItemRowCol(y, x int) terrain {
	pt := xy{x, y}
	return g[pt]
}

type xy struct {
	x, y int
}

type terrain struct {
	ch rune
}

func (t terrain) String() string {
	return fmt.Sprintf("%c", t.ch)
}

func print(g grid) {
	mx, my := maxes(g)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			fmt.Print(string(g[xy{x, y}].ch))
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 13 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	gs := parseInput(in)
	return sumReflections(gs)
}

func part2(in input) output {
	return 0
}

func sumReflections(gs []grid) int {
	total := 0
	for _, g := range gs {
		total += getReflections(g)
	}
	return total
}

func getReflections(m grid) int {
	mx, my := maxes(m)
	fmt.Println("---- check cols------")
	midx := getMirrorMidpoint(m, mx, my, m.ItemColRow)
	fmt.Println("---- check rows------")
	midy := getMirrorMidpoint(m, my, mx, m.ItemRowCol)
	return 100*midy + midx
}

func getMirrorMidpoint(g grid, max1, max2 int, getitem func(i, j int) terrain) int {
	m := make(map[int]int)
	mid := 0
	cur := max1
	for i := 0; i <= max1; i++ {
		for i2 := cur; i2 > i; i2-- {
			diffs := 0
			for j := 0; j <= max2; j++ {
				item1 := getitem(i, j)
				item2 := getitem(i2, j)

				if item1.ch != item2.ch {
					diffs++
				}
			}

			if diffs == 0 {
				m[i] = i2
				m[i2] = i
				_, haszero := m[0]
				_, hasmax := m[max1]
				fmt.Println(cur, i, i2, haszero, hasmax, m)
				if cur-i2 == 1 && (haszero || hasmax) {
					fmt.Println("got thing", i)
					mid = i + 1
					fmt.Println("newmid", mid)
				}
				cur--
				break
			}
		}
	}

	print(g)
	fmt.Println(m)
	fmt.Println("mid", mid)
	return mid
}

func maxes(m grid) (int, int) {
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

func parseInput(in input) []grid {
	grids := []grid{}
	cur := make(grid)
	y := 0
	for i, line := range in {
		if line == "" {
			grids = append(grids, cur)
			cur = make(grid)
			y = 0
			continue
		}

		for x, c := range line {
			pt := xy{x, y}
			t := terrain{ch: c}
			cur[pt] = t
		}
		y++
		if i == len(in)-1 {
			grids = append(grids, cur)
		}
	}
	return grids
}
