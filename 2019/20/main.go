package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int
type grid [][]block

type blockType string

const (
	wall  blockType = "#"
	blank blockType = " "
	path  blockType = "."
)

type portalType string

const (
	inner portalType = "inner"
	outer portalType = "outer"
)

type portal struct {
	id    string
	pt    xy
	ptype portalType
}

func (p portal) String() string {
	return fmt.Sprintf("id: %s xy: %v type: %s", p.id, p.pt, p.ptype)
}

type block struct {
	pt     xy
	t      blockType
	portal *portal
}

type xy struct {
	x, y int
}

type xylevel struct {
	xy
	level int
}

func (p xy) String() string {
	return fmt.Sprintf("(x: %d, y: %d)", p.x, p.y)
}

type state struct {
	moves []xy
	pt    xy
	goal  bool
	level int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	g := parseGrid(in)
	portals := getPortals(g)

	shortest := traverse(g, portals["AA"][0].pt, portals["ZZ"][0].pt, portals, false)
	return shortest
}

func part2(in input) output {
	g := parseGrid(in)
	portals := getPortals(g)
	shortest := traverse(g, portals["AA"][0].pt, portals["ZZ"][0].pt, portals, true)
	return shortest
}

func traverse(g grid, start, end xy, portals map[string][]portal, recurse bool) int {
	goalLevel := 1
	states := []*state{&state{pt: start, level: 1}}
	goals := []*state{}
	visited := make(map[xylevel]bool)
	minpath := 100000

	for len(states) > 0 {
		cur := states[0]
		states = states[1:]

		if len(cur.moves) > minpath {
			continue
		}

		mv := cur.pt
		b := g[mv.y][mv.x]

		leveldelta := 0
		mvs := getValidMoves(g, b.pt)
		var exit portal
		if b.portal != nil && b.portal.pt != start {
			add := true
			if recurse {
				if cur.level == goalLevel && b.portal.ptype == outer {
					add = false
				} else if cur.level > goalLevel && (b.portal.pt == start || b.portal.pt == end) {
					add = false
				}
			}

			if add {
				plist := portals[b.portal.id]
				for _, p := range plist {
					if p.pt == b.portal.pt {
						continue
					}

					exit = p
				}

				if recurse {
					if b.portal.ptype == inner {
						leveldelta = 1
					} else {
						leveldelta = -1
					}
				}

				mvs = append(mvs, exit.pt)
			}
		}

		for _, mv := range mvs {
			s := copyState(cur)
			xlv := xylevel{xy: mv, level: s.level}
			if recurse && xlv.xy == exit.pt {
				xlv.level += leveldelta
			}

			if _, ok := visited[xlv]; ok {
				continue
			}

			visited[xlv] = true
			s.pt = mv
			s.moves = append(s.moves, mv)
			s.level = xlv.level

			if mv == end && s.level == goalLevel {
				s.goal = true
				if len(s.moves) < minpath {
					minpath = len(s.moves)
				}
				goals = append(goals, s)
			} else {
				states = append(states, s)
			}
		}
	}

	return minpath
}

func copyState(st *state) *state {
	cp := &state{level: st.level}
	cp.pt = st.pt
	for _, mv := range st.moves {
		cp.moves = append(cp.moves, mv)
	}
	return cp
}

func getValidMoves(g grid, pt xy) []xy {
	mvs := []xy{
		{pt.x + 1, pt.y}, {pt.x - 1, pt.y}, {pt.x, pt.y + 1}, {pt.x, pt.y - 1},
	}

	valid := []xy{}
	for _, mv := range mvs {
		if mv.x < 0 || mv.y < 0 || mv.x > len(g[0])-1 || mv.y > len(g)-1 {
			continue
		}

		b := g[mv.y][mv.x]
		if b.t == blank || b.t == wall {
			continue
		}

		valid = append(valid, mv)
	}
	return valid
}

func getPortals(g grid) map[string][]portal {
	portals := make(map[string][]portal)
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			b := g[y][x]

			if b.portal != nil {
				portals[b.portal.id] = append(portals[b.portal.id], *b.portal)
			}
		}
	}
	return portals
}

func parseGrid(in input) grid {
	g := make(grid, len(in))
	lreg := regexp.MustCompile("[a-zA-Z]{2}")

	for y, line := range in {
		g[y] = make([]block, len(line))
		for x, c := range line {
			b := block{pt: xy{x, y}}
			switch c {
			case '#':
				b.t = wall
			case '.':
				b.t = path
			default:
				b.t = blank
			}

			var prevx, nextx string
			var prevy, nexty string
			if b.t == path {
				if x > 1 {
					prevx = string([]byte{line[x-2], line[x-1]})
				}

				if x < len(line)-2 {
					nextx = string([]byte{line[x+1], line[x+2]})
				}

				if y > 0 {
					prevy = string([]byte{in[y-2][x], in[y-1][x]})
				}

				if y < len(in)-2 {
					nexty = string([]byte{in[y+1][x], in[y+2][x]})
				}
			}

			for _, s := range []string{prevx, nextx, prevy, nexty} {
				if lreg.MatchString(s) {
					ptype := inner

					if x < 3 || x > len(line)-4 || y < 3 || y > len(g)-4 {
						ptype = outer
					}
					b.portal = &portal{id: s, pt: b.pt, ptype: ptype}
				}
			}

			g[y][x] = b
		}
	}

	return g
}
