package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

var input = "24.txt"

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

	lines := []string{}

	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()
			lines = append(lines, txt)
		}
	}

	ducts := make([][]Duct, len(lines))

	startx, starty := 0, 0
	targets := 0
	goals := []Point{}

	for i := 0; i < len(lines); i++ {
		ducts[i] = make([]Duct, len(lines[i]))

		for j := 0; j < len(lines[i]); j++ {
			duct := Duct{C: rune(lines[i][j]), Open: lines[i][j] != '#', Target: lines[i][j] != '#' && lines[i][j] != '.'}
			ducts[i][j] = duct

			if ducts[i][j].C == '0' {
				startx, starty = j, i
			}
			if ducts[i][j].Target {
				targets++
				if ducts[i][j].C != '0' {
					goals = append(goals, Point{X: j, Y: i})
				}
			}
		}
	}

	print(ducts)

	var perms [][]Point = Permutate(goals)
	start := Point{X: startx, Y: starty}

	part2 := true

	lengths := []int{}

	for _, list := range perms {
		alllist := append([]Point{start}, list...)
		if part2 {
			alllist = append(alllist, start)
		}
		solvelen := 0
		for i := 1; i < len(alllist); i++ {
			visited := make(map[Point]bool)
			solves := solve(alllist[i-1], alllist[i], ducts, visited)

			pathlens := []int{}
			for _, ln := range solves {
				pathlens = append(pathlens, ln.Moves)

			}

			sort.Ints(pathlens)
			solvelen += pathlens[0]
		}

		lengths = append(lengths, solvelen)
	}

	sort.Ints(lengths)
	fmt.Println(lengths[0])

	fmt.Println("Time", time.Since(startTime))
}

func print(ducts [][]Duct) {
	for y := 0; y < len(ducts); y++ {
		for x := 0; x < len(ducts[y]); x++ {
			ch := '#'
			if ducts[y][x].Open {
				ch = ducts[y][x].C
			}
			fmt.Print(string(ch))
		}
		fmt.Print("\n")
	}
}

func solve(point, goal Point, ducts [][]Duct, visited map[Point]bool) []State {
	queue := []State{}
	queue = append(queue, State{Moves: 0, Point: point})
	solves := []State{}
	minsolve := 10000

	for len(queue) > 0 {
		state := queue[0]
		queue = queue[1:]

		moves := getMoves(state.Point, ducts)

		for _, mv := range moves {
			mvstate := State{Moves: state.Moves + 1, Point: mv}
			if mv.X == goal.X && mv.Y == goal.Y {
				minsolve = int(math.Min(float64(minsolve), float64(state.Moves)))
				solves = append(solves, mvstate)
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

	return solves
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
