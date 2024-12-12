package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

type plot struct {
	c     rune
	edges int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 12 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parse(in)
	regions := findRegions(m)
	return getTotalFences(regions)
}

func part2(in input) output {
	return 0
}

func getTotalFences(regions []map[xy]plot) int {
	total := 0
	for _, r := range regions {
		total += len(r) * getRegionPerimeter(r)
	}
	return total
}

func getRegionPerimeter(m map[xy]plot) int {
	perimeter := 0
	for cur := range m {
		sides := []xy{{cur.x + 1, cur.y}, {cur.x - 1, cur.y}, {cur.x, cur.y + 1}, {cur.x, cur.y - 1}}
		total := 4
		for _, s := range sides {
			if _, ok := m[s]; ok {
				total--
			}
		}
		perimeter += total
	}
	return perimeter
}

func findRegions(m map[xy]plot) []map[xy]plot {
	regions := []map[xy]plot{}

	tested := make(map[xy]bool)

	for cur := range m {
		if _, ok := tested[cur]; ok {
			continue
		}

		region := getRegionAt(m, cur)
		for k := range region {
			tested[k] = true
		}
		regions = append(regions, region)
	}

	return regions
}

func getRegionAt(m map[xy]plot, pt xy) map[xy]plot {
	region := make(map[xy]plot)
	visited := make(map[xy]bool)

	c := m[pt].c

	queue := []xy{pt}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if _, ok := visited[cur]; ok {
			continue
		}

		if b, ok := m[cur]; !ok || b.c != c {
			continue
		}

		visited[cur] = true
		region[cur] = m[cur]

		mvs := []xy{{cur.x + 1, cur.y}, {cur.x - 1, cur.y}, {cur.x, cur.y + 1}, {cur.x, cur.y - 1}}
		queue = append(queue, mvs...)
	}

	return region
}

func parse(in []string) map[xy]plot {
	m := make(map[xy]plot)
	for y := 0; y < len(in); y++ {
		for x := 0; x < len(in[y]); x++ {
			pt := xy{x, y}
			m[pt] = plot{c: rune(in[y][x])}
		}
	}
	return m
}
