package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

var input = "10.txt"

type contents int

const (
	blank    contents = 0
	asteroid contents = 1
)

type visibility int

const (
	invisible visibility = 0
	visible   visibility = 1
	self      visibility = 2
)

type point struct {
	x, y    int
	content contents
}

type result struct {
	point      point
	visibility visibility
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d) %s", p.x, p.y, p.content)
}

func (c contents) String() string {
	s := "blank"
	if c == asteroid {
		s = "asteroid"
	}
	return s
}

func (p point) awayvector(p2, max point) point {
	v := point{x: p2.x - p.x, y: p2.y - p.y}
	xs := math.Signbit(float64(v.x))
	ys := math.Signbit(float64(v.y))
	absx, absy := int(math.Abs(float64(v.x))), int(math.Abs(float64(v.y)))
	if v.x == 0 || v.y == 0 {
		if v.x != 0 {
			v.x = 1
			if xs {
				v.x = -1
			}
		} else if v.y != 0 {
			v.y = 1
			if ys {
				v.y = -1
			}
		}
	} else if absx != absy {
		for d := max.x; d > 1; d-- {
			if absx % d == 0 && absy % d == 0 {
				absx = absx / d
				absy = absy / d

				v.x = absx
				if xs {
					v.x = -v.x
				}
	
				v.y = absy
				if ys {
					v.y = -v.y
				}
			}
	
		}

	} else if absx == absy {
		v.x = 1
		v.y = 1
		if xs {
			v.x = -1
		}
		if ys {
			v.y = -1
		}
	}
	return v
}

func (p point) nextAway(v point) point {
	return point{x: p.x + v.x, y: p.y + v.y}
}

func (p point) inGraph(maxp point) bool {
	return p.x >= 0 && p.y >= 0 && p.x <= maxp.x && p.y <= maxp.y
}

func (r result) String() string {
	return fmt.Sprintf("%v %v", r.point, r.visibility)
}

func (v visibility) String() string {
	if v == visible {
		return "visible"
	} else if v == invisible {
		return "invisible"
	} else {
		return "self"
	}
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	space := [][]point{}
	for r, line := range lines {
		space = append(space, []point{})
		sline := space[len(space)-1]
		for c, ch := range line {
			cnt := blank
			if ch == '#' {
				cnt = asteroid
			}
			p := point{x: c, y: r, content: cnt}
			sline = append(sline, p)
		}
		space[len(space)-1] = sline
	}

	p, count := getResults(space)

	fmt.Println("Part 1: ", count, p)
	fmt.Println("Time", time.Since(startTime))
}

func getResults(space [][]point) (point, int) {
	max := 0
	var best point
	for r, line := range space {
		for c := range line {
			if space[r][c].content == asteroid {
				p := point{x: c, y: r}
				vis, _ := countVisible(space, p)

				if vis > max {
					best = p
					max = vis
				}
			}
		}
	}

	return best, max
}

func countVisible(space [][]point, p point) (int, [][]result) {
	maxx, maxy := len(space[0])-1, len(space)-1
	maxp := point{x: maxx, y: maxy}
	results := make([][]result, len(space))
	for i := 0; i < len(results); i++ {
		results[i] = make([]result, len(space[i]))

		for j := 0; j < len(results[i]); j++ {
			results[i][j].point = space[i][j]

			if space[i][j].content == asteroid {
				results[i][j].visibility = visible
				if p.y == i && p.x == j {
					results[i][j].visibility = self
				}
			}
		}
	}

	for r1 := 0; r1 < len(space); r1++ {
		for c1 := 0; c1 < len(space[r1]); c1++ {
			if space[r1][c1].content == asteroid && results[r1][c1].visibility == visible {
				p2 := point{x: c1, y: r1}

				v := p.awayvector(p2, maxp)

				next := p2.nextAway(v)
				for next.inGraph(maxp) {
					results[next.y][next.x].visibility = invisible
					next = next.nextAway(v)
				}
			}
		}
	}

	count := 0
	for i := range space {
		for j := 0; j < len(results[i]); j++ {
			if results[i][j].visibility == visible && space[i][j].content == asteroid {
				count++
			}
		}
	}
	return count, results
}
