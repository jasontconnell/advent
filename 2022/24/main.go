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

var (
	up    xy = xy{0, -1}
	down  xy = xy{0, 1}
	left  xy = xy{-1, 0}
	right xy = xy{1, 0}
)

type xy struct {
	x, y int
}

type block struct {
	wall bool
}

type blizzard struct {
	dir xy
}

type blizzards map[xy][]blizzard

type state struct {
	pt     xy
	minute int
}

type move struct {
	pt xy
}

type cachekey struct {
	minute int
	pt     xy
}

func print(grid map[xy]block, bzs blizzards) {
	min, max := minmax(grid)

	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			pt := xy{x, y}
			c := "."
			if g, ok := grid[pt]; ok && g.wall {
				c = "#"
			}
			if list, ok := bzs[pt]; ok {
				if len(list) == 1 {
					switch list[0].dir {
					case down:
						c = "v"
					case up:
						c = "^"
					case left:
						c = "<"
					case right:
						c = ">"
					}
				} else {
					c = strconv.Itoa(len(list))
				}
			}
			fmt.Print(c)
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
	fmt.Fprintln(w, "--2022 day 24 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	g, b := parseInput(in)
	start, end := getStartEnd(g)
	bstates := generateBlizzardStates(g, b, 500)
	minutes, _ := simulate(0, g, start, end, bstates)
	return minutes
}

func part2(in input) output {
	g, b := parseInput(in)
	start, end := getStartEnd(g)
	minutes := 0
	bstates := generateBlizzardStates(g, b, 1100)
	for i := 0; i < 3; i++ {
		minutes, b = simulate(minutes, g, start, end, bstates)
		start, end = end, start
	}
	return minutes
}

func generateBlizzardStates(grid map[xy]block, blz blizzards, num int) []blizzards {
	// generate states for num blizzard moves
	// instead of each branch of the search
	bstates := []blizzards{blz}
	for i := 1; i < num; i++ {
		bstates = append(bstates, moveBlizzards(grid, bstates[i-1]))
	}
	return bstates
}

func getStartEnd(grid map[xy]block) (xy, xy) {
	min, max := minmax(grid)
	var start, end xy
	for x := min.x; x <= max.x; x++ {
		if b, ok := grid[xy{x, min.y}]; ok && !b.wall {
			start = xy{x, min.y}
		}
		if b, ok := grid[xy{x, max.y}]; ok && !b.wall {
			end = xy{x, max.y}
		}
	}
	return start, end
}

func simulate(minute int, grid map[xy]block, start, end xy, bstates []blizzards) (int, blizzards) {
	visit := make(map[cachekey]bool)
	queue := common.NewPriorityQueue(func(s state) float64 {
		return 1/dist(s.pt, end) + 1/float64(s.minute)
	})

	queue.Enqueue(state{minute: minute, pt: start})
	var best state = state{minute: math.MaxInt32}
	for queue.Any() {
		cur := queue.Dequeue()

		key := cachekey{minute: cur.minute, pt: cur.pt}
		if _, ok := visit[key]; ok {
			continue
		}
		visit[key] = true

		if cur.pt == end && cur.minute < best.minute {
			best = cur
		} else if cur.pt == end || cur.minute >= best.minute {
			continue
		}

		if cur.minute > len(bstates)-2 {
			continue
		}

		if list, ok := bstates[cur.minute][cur.pt]; ok && len(list) > 0 {
			continue
		}

		mvs := getMoves(cur.pt, grid, bstates[cur.minute+1])
		for _, mv := range mvs {
			st := state{pt: mv.pt, minute: cur.minute + 1}
			queue.Enqueue(st)
		}
	}
	return best.minute, bstates[best.minute]
}

func getMoves(pt xy, grid map[xy]block, blz blizzards) []move {
	mvs := []move{}
	for _, d := range []xy{up, down, left, right} {
		np := xy{pt.x + d.x, pt.y + d.y}
		if bk, ok := grid[np]; ok && !bk.wall {
			if list, ok := blz[np]; !ok || len(list) == 0 {
				mvs = append(mvs, move{pt: np})
			}
		}
	}
	mvs = append(mvs, move{pt: pt})
	return mvs
}

func dist(p1, p2 xy) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return math.Abs(float64(dx)) + math.Abs(float64(dy))
}

func moveBlizzards(grid map[xy]block, blz blizzards) blizzards {
	min, max := minmax(grid)
	cp := make(blizzards)
	for k, v := range blz {
		for _, b := range v {
			np := xy{k.x + b.dir.x, k.y + b.dir.y}
			if bk, ok := grid[np]; ok && bk.wall {
				if b.dir.x != 0 {
					if np.x == max.x {
						np.x = min.x + 1
					} else {
						np.x = max.x - 1
					}
				} else {
					if np.y == max.y {
						np.y = min.y + 1
					} else {
						np.y = max.y - 1
					}
				}
			}
			cp[np] = append(cp[np], b)
		}
	}
	return cp
}

func minmax[T any](m map[xy]T) (xy, xy) {
	min, max := xy{math.MaxInt32, math.MaxInt32}, xy{math.MinInt32, math.MinInt32}
	for k := range m {
		if k.x < min.x {
			min.x = k.x
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.y > max.y {
			max.y = k.y
		}
	}
	return min, max
}

func parseInput(in input) (map[xy]block, blizzards) {
	grid := make(map[xy]block)
	bzs := make(blizzards)
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			b := block{wall: c == '#'}
			grid[pt] = b

			if c != '.' && c != '#' {
				var d xy
				switch c {
				case 'v':
					d = down
				case '^':
					d = up
				case '<':
					d = left
				case '>':
					d = right
				}

				bz := blizzard{dir: d}
				bzs[pt] = []blizzard{bz}
			}
		}
	}
	return grid, bzs
}
