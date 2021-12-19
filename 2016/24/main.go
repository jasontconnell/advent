package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Duct struct {
	Open   bool
	Target bool
	C      rune
}

func (duct Duct) String() string {
	return string(duct.C)
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

type State struct {
	Point Point
	Moves int
}

func (s State) String() string {
	return fmt.Sprintf("Reached %v in %v moves", s.Point, s.Moves)
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 24 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	ducts, goals, start := parseInput(in)
	return run(start, goals, ducts, false)
}

func part2(in input) output {
	ducts, goals, start := parseInput(in)
	return run(start, goals, ducts, true)
}

func run(start Point, goals []Point, ducts [][]Duct, returnToStart bool) int {
	var perms [][]Point = Permutate(goals)

	lengths := []int{}

	for _, list := range perms {
		alllist := append([]Point{start}, list...)
		if returnToStart {
			alllist = append(alllist, start)
		}
		solvelen := 0
		for i := 1; i < len(alllist); i++ {
			visited := make(map[Point]bool)
			minsolve := solve(alllist[i-1], alllist[i], ducts, visited)

			solvelen += minsolve
		}

		lengths = append(lengths, solvelen)
	}

	sort.Ints(lengths)
	return lengths[0]
}

func solve(point, goal Point, ducts [][]Duct, visited map[Point]bool) int {
	queue := []State{}
	queue = append(queue, State{Moves: 0, Point: point})
	minsolve := math.MaxInt32

	for len(queue) > 0 {
		state := queue[0]
		queue = queue[1:]

		moves := getMoves(state.Point, ducts)

		for _, mv := range moves {
			mvstate := State{Moves: state.Moves + 1, Point: mv}
			if mv == goal {
				if mvstate.Moves < minsolve {
					minsolve = mvstate.Moves
				}
				continue
			}

			if _, exists := visited[mv]; !exists {
				visited[mv] = true

				if state.Moves+1 < minsolve {
					queue = append(queue, mvstate)
				}
			}
		}
	}

	return minsolve
}

func contains(rs []rune, r rune) bool {
	c := false
	for _, s := range rs {
		c = c || s == r
	}
	return c
}

func getMoves(point Point, ducts [][]Duct) []Point {
	pts := []Point{}
	for _, pt := range []Point{Point{X: point.X, Y: point.Y + 1}, Point{X: point.X, Y: point.Y - 1}, Point{X: point.X + 1, Y: point.Y}, Point{X: point.X - 1, Y: point.Y}} {
		if ducts[pt.Y][pt.X].Open {
			pts = append(pts, pt)
		}
	}
	return pts
}

func Permutate(str []Point) [][]Point {
	var ret [][]Point

	if len(str) == 2 {
		ret = append(ret, []Point{str[0], str[1]})
		ret = append(ret, []Point{str[1], str[0]})
	} else {

		for i := 0; i < len(str); i++ {
			strc := make([]Point, len(str))
			copy(strc, str)

			t := strc[i]
			sh := append(strc[:i], strc[i+1:]...)
			perm := Permutate(sh)

			for _, p := range perm {
				p = append([]Point{t}, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}

func parseInput(in input) ([][]Duct, []Point, Point) {
	ducts := make([][]Duct, len(in))

	targets := 0
	goals := []Point{}
	start := Point{}

	for i := 0; i < len(in); i++ {
		ducts[i] = make([]Duct, len(in[i]))

		for j := 0; j < len(in[i]); j++ {
			duct := Duct{C: rune(in[i][j]), Open: in[i][j] != '#', Target: in[i][j] != '#' && in[i][j] != '.'}
			ducts[i][j] = duct

			if ducts[i][j].C == '0' {
				start.X = j
				start.Y = i
			}
			if ducts[i][j].Target {
				targets++
				if ducts[i][j].C != '0' {
					goals = append(goals, Point{X: j, Y: i})
				}
			}
		}
	}

	return ducts, goals, start
}
