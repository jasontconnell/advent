package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 23 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	machines := parse(in)
	return findTs(machines)
}

func part2(in input) output {
	return 0
}

func findTs(g common.Graph[string, int]) int {
	ts := 0
	vertices := g.GetVertices()
	sort.Strings(vertices)
	threes := [][]string{}
	dedup := make(map[string]bool)
	for _, v := range vertices {

		neighbors := g.Neighbors(v)
		sort.Strings(neighbors)
		for _, n := range neighbors {
			n2 := g.Neighbors(n)
			sort.Strings(n2)
			for _, n2n := range n2 {
				if g.Adjacent(v, n2n) {
					three := []string{v, n, n2n}
					sort.Strings(three)
					s := three[0] + three[1] + three[2]
					if _, ok := dedup[s]; ok {
						continue
					}
					dedup[s] = true
					threes = append(threes, three)
				}
			}
		}
	}
	sort.Slice(threes, func(i, j int) bool {
		return threes[i][0] < threes[j][0]
	})
	for _, t := range threes {
		for _, s := range t {
			if strings.HasPrefix(s, "t") {
				ts++
				break
			}
		}
	}

	return ts
}

func findCliques(g common.Graph[string, int]) {
	// queue := g.GetVertices()

}

func parse(in []string) common.Graph[string, int] {
	g := common.NewGraph[string]()

	for _, line := range in {
		ms := strings.Split(line, "-")
		g.AddVertices(ms[0], ms[1])
		g.AddEdge(ms[0], ms[1])
	}

	return g
}
