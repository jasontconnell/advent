package main

import (
	"fmt"
	"log"
	"os"

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
	fmt.Fprintln(w, "--2024 day 04 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	return findOccurrences("XMAS", parse(in))
}

func part2(in input) output {
	total := findDiagonals("MAS", 'A', parse(in))
	return total
}

func findOccurrences(word string, graph map[xy]byte) int {
	count := 0
	for k, v := range graph {
		if v != word[0] {
			continue
		}
		ss := getWordsAt(k, graph, len(word))
		for _, s := range ss {
			if s == word {
				count++
			}
		}
	}
	return count
}

type xy struct {
	x, y int
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

func findDiagonals(word string, key byte, graph map[xy]byte) int {
	total := 0
	deltas := [][]xy{
		{{1, 1}, {0, 0}, {-1, -1}},
		{{-1, -1}, {0, 0}, {1, 1}},
		{{1, -1}, {0, 0}, {-1, 1}},
		{{-1, 1}, {0, 0}, {1, -1}},
	}
	for k, _ := range graph {
		if graph[k] != key {
			continue
		}
		count := 0

		for _, da := range deltas {
			w := ""
			for i := 0; i < len(da); i++ {
				c := k.add(da[i])
				if b, ok := graph[c]; ok {
					w += string(b)
				}
			}
			if w == word {
				count++
			}
		}

		if count == 2 {
			total++
		}
	}

	return total
}

func getWordsAt(pt xy, graph map[xy]byte, maxlen int) []string {
	words := []string{}
	deltas := []xy{{1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}, {0, 1}, {0, -1}}
	cur := pt
	for _, d := range deltas {
		w := ""
		add := true
		for i := 0; i < maxlen; i++ {
			if c, ok := graph[cur]; ok {
				w += string(c)
			} else {
				add = false
			}
			cur = cur.add(d)
		}
		if add {
			words = append(words, w)
		}
		cur = pt
	}
	return words
}

func parse(graph []string) map[xy]byte {
	m := make(map[xy]byte)
	for y := 0; y < len(graph); y++ {
		for x := 0; x < len(graph[y]); x++ {
			m[xy{x, y}] = graph[y][x]
		}
	}
	return m
}
