package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xyz struct {
	x, y, z int
}

type pair struct {
	p1, p2     xyz
	dist       float64
	idx1, idx2 int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 08 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	process := 1000
	if len(in) < 100 {
		process = 10
	}
	coords := parseInput(in)
	groups := groupPairs(coords, process)

	for _, g := range groups {
		log.Println("group", len(g))
	}

	size := 1
	for i := len(groups) - 1; i > len(groups)-4; i-- {
		size *= len(groups[i])
	}
	return size
}

func part2(in input) output {
	return 0
}

func groupPairs(coords []xyz, process int) []map[xyz]bool {
	groups := []map[xyz]bool{}

	cp := make([]xyz, len(coords))
	copy(cp, coords)

	pairs := getDistances(cp)

	for _, next := range pairs[:process+1] {
		var group map[xyz]bool
		for _, g := range groups {
			if _, ok := g[next.p1]; ok {
				group = g
				break
			} else if _, ok := g[next.p2]; ok {
				group = g
				break
			}
		}

		if group == nil {
			group = make(map[xyz]bool)
			groups = append(groups, group)
		}

		if group[next.p1] && group[next.p2] {
			continue
		}

		group[next.p1] = true
		group[next.p2] = true
	}

	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) < len(groups[j])
	})

	return groups
}

func getDistances(coords []xyz) []pair {
	pairs := []pair{}
	visit := make(map[xyz]bool)
	for i := 0; i < len(coords); i++ {
		for j := 0; j < len(coords); j++ {
			if i == j {
				continue
			}
			if _, ok := visit[xyz{i, j, 0}]; ok {
				continue
			}
			if _, ok := visit[xyz{j, i, 0}]; ok {
				continue
			}
			visit[xyz{i, j, 0}] = true

			p1, p2 := coords[i], coords[j]
			dist := distance(p1, p2)

			pairs = append(pairs, pair{p1: p1, p2: p2, dist: dist, idx1: i, idx2: j})
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].dist < pairs[j].dist
	})

	return pairs
}

func distance(p1, p2 xyz) float64 {
	xd := p1.x - p2.x
	yd := p1.y - p2.y
	zd := p1.z - p2.z
	return math.Sqrt(float64(xd*xd + yd*yd + zd*zd))
}

func parseInput(in input) []xyz {
	coords := []xyz{}
	for _, line := range in {
		sp := strings.Split(line, ",")
		x, _ := strconv.Atoi(sp[0])
		y, _ := strconv.Atoi(sp[1])
		z, _ := strconv.Atoi(sp[2])
		coords = append(coords, xyz{x, y, z})
	}
	return coords
}
