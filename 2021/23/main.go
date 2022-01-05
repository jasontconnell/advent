package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type pqueue []state

func (pq pqueue) Len() int { return len(pq) }

func (pq pqueue) Less(i, j int) bool {
	return pq[i].score < pq[j].score
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	st := old[n-1]
	*pq = old[0 : n-1]
	return st
}

func (pq *pqueue) Push(x interface{}) {
	st := x.(state)
	*pq = append(*pq, st)
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

const hallwayy int = 1

type state struct {
	energyuse int
	g         grid
	score     int
	rooms     *map[podname]room
	pods      map[xy]pod
	loc       xy
	moves     []move
}

type xy struct {
	x, y int
}

type podname string

const (
	None podname = "N"
	A    podname = "A"
	B    podname = "B"
	C    podname = "C"
	D    podname = "D"
)

type pod struct {
	name      podname
	energyuse int
	num       int
}

type visitkey struct {
	name podname
	num  int
	pt   xy
}

var Empty pod = pod{name: None, energyuse: 0, num: 0}

type room struct {
	owner podname
	pts   map[xy]bool
	x     int
	y     int
}

type block struct {
	room    bool
	hallway bool
	wall    bool
	pod     pod
}

type grid [][]block

type move struct {
	pod   pod
	pos   xy
	to    xy
	steps int
}

func printGrid(g grid) {
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			b := g[y][x]

			s := " "
			if b.wall {
				s = "#"
			} else if b.hallway || b.room {
				s = "."
			}

			if b.pod.name != None {
				s = string(b.pod.name)
			}
			fmt.Print(s)
		}
		fmt.Println()
	}
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
	fmt.Fprintln(w, "--2021 day 23 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	g, pods, rooms := parseInput(in)
	return solve(g, pods, rooms)
}

func part2(in input) output {
	in2 := make(input, len(in))
	copy(in2, in[:3])
	in2[3] = "  #D#C#B#A#"
	in2[4] = "  #D#B#A#C#"
	in2 = append(in2, in[3:]...)
	g, pods, rooms := parseInput(in2)
	return solve(g, pods, rooms)
}

func allMovablePods(g grid, rooms map[podname]room, pods map[xy]pod) []xy {
	candidates := []xy{}
	for pt := range pods {
		if isFinal(g, rooms, pt) {
			continue
		}
		b := g[pt.y][pt.x]
		bup := g[pt.y-1][pt.x]

		if g[hallwayy][pt.x-1].pod != Empty && g[hallwayy][pt.x+1].pod != Empty {
			continue
		}

		if b.room && bup.pod == Empty {
			candidates = append(candidates, pt)
		}
	}
	return candidates
}

func roomablePods(g grid, rooms map[podname]room, pods map[xy]pod) []xy {
	candidates := []xy{}
	for pt := range pods {
		b := g[pt.y][pt.x]
		if b.hallway {
			candidates = append(candidates, pt)
		}
	}
	roomable := []xy{}
	for _, pt := range candidates {
		b := g[pt.y][pt.x]

		if !b.hallway {
			continue
		}

		r := rooms[b.pod.name]

		roomavail := true
		var empty xy
		for rpt := range r.pts {
			rb := g[rpt.y][rpt.x]
			if rb.pod.name != b.pod.name && rb.pod != Empty {
				roomavail = false
				break
			}
			if rb.pod == Empty && rpt.y > empty.y {
				empty = rpt
			}
		}

		if !roomavail {
			continue
		}

		if pathClear(g, pt, empty) {
			roomable = append(roomable, pt)
		}
	}
	return roomable
}

func getKey(g grid) string {
	s := ""
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			b := g[y][x]
			if b.hallway || b.room {
				s += string(b.pod.name)
			}
		}
	}
	return s
}

func solve(g grid, pods map[xy]pod, rooms map[podname]room) int {
	gcp := make(grid, len(g))
	for y := range gcp {
		gcp[y] = make([]block, len(g[y]))
		copy(gcp[y], g[y])
	}

	queue := pqueue{}

	movable := allMovablePods(g, rooms, pods)

	for _, pt := range movable {
		initial := state{g: g, energyuse: 0, pods: pods, rooms: &rooms, loc: pt, moves: []move{}}
		queue.Push(initial)
	}
	heap.Init(&queue)

	var solvestate state
	visited := map[string]int{}
	minenergy := math.MaxInt32
	for queue.Len() > 0 {
		cur := queue.Pop().(state)

		key := getKey(cur.g)
		if n, ok := visited[key]; ok && cur.energyuse > n || cur.energyuse > minenergy {
			continue
		}
		visited[key] = cur.energyuse

		roomable := roomablePods(cur.g, *cur.rooms, cur.pods)
		for len(roomable) != 0 {
			total := 0
			for _, pt := range roomable {
				mv := getRoomMove(cur.g, *cur.rooms, pt)
				applyMove(&cur, mv)
				pd := cur.pods[mv.pos]
				total += pd.energyuse * mv.steps
			}
			roomable = roomablePods(cur.g, *cur.rooms, cur.pods)
			cur.score += total
		}

		if isSolved(cur.g, rooms) {
			if cur.energyuse < minenergy {
				minenergy = cur.energyuse
				solvestate = cur
			}
			continue
		}

		mvs := getMoves(cur.g, cur.loc, rooms)
		for _, mv := range mvs {
			st := copyState(cur)
			applyMove(&st, mv)
			queue.Push(st)
		}

		if len(mvs) == 0 {
			mpods := allMovablePods(cur.g, *cur.rooms, cur.pods)
			for _, pod := range mpods {
				st2 := copyState(cur)
				st2.loc = pod
				queue.Push(st2)
			}
		}

	}

	fmt.Println("----- winning moves -----")
	printGrid(gcp)
	for _, mv := range solvestate.moves {
		fmt.Println()
		mp := gcp[mv.pos.y][mv.pos.x].pod
		gcp[mv.pos.y][mv.pos.x].pod = Empty
		gcp[mv.to.y][mv.to.x].pod = mp
		printGrid(gcp)
	}

	return minenergy
}

func applyMove(st *state, mv move) {
	b := st.g[mv.pos.y][mv.pos.x]
	energy := b.pod.energyuse
	st.energyuse += energy * mv.steps
	st.g[mv.pos.y][mv.pos.x].pod = Empty
	st.g[mv.to.y][mv.to.x].pod = b.pod
	st.loc = mv.to
	delete(st.pods, mv.pos)
	st.pods[st.loc] = b.pod
	st.moves = append(st.moves, mv)
	st.score = energy * mv.steps
}

func copyState(cur state) state {
	st := state{energyuse: cur.energyuse, rooms: cur.rooms, loc: cur.loc, score: cur.score}

	st.g = make(grid, len(cur.g))
	for y, row := range cur.g {
		st.g[y] = make([]block, len(row))
		copy(st.g[y], cur.g[y])
	}

	st.pods = map[xy]pod{}
	for k, v := range cur.pods {
		st.pods[k] = v
	}

	st.moves = make([]move, len(cur.moves))
	copy(st.moves, cur.moves)
	return st
}

func dist(p1, p2 xy) int {
	d := math.Abs(float64(p1.x-p2.x)) + math.Abs(float64(p1.y-p2.y))
	return int(d)
}

func getMoves(g grid, cur xy, rooms map[podname]room) []move {
	mvs := []move{}

	b := g[cur.y][cur.x]
	pmvs := getPodMoves(g, rooms, cur)
	for _, m := range pmvs {
		steps := dist(cur, m)
		mvs = append(mvs, move{pod: b.pod, pos: cur, to: m, steps: steps})
	}

	return mvs
}

func isFinal(g grid, rooms map[podname]room, pos xy) bool {
	b := g[pos.y][pos.x]
	r := rooms[b.pod.name]
	_, ok := r.pts[pos]
	if !ok {
		return false
	}

	for p := range r.pts {
		rb := g[p.y][p.x]

		if rb.pod != Empty && rb.pod.name != b.pod.name {
			return false
		}
	}

	return true
}

func getRoomMove(g grid, rooms map[podname]room, pos xy) move {
	b := g[pos.y][pos.x]
	var mv move
	if !b.hallway {
		return mv
	}

	r := rooms[b.pod.name]

	var empty xy
	for pt := range r.pts {
		bp := g[pt.y][pt.x]
		if bp.pod == Empty && pt.y > empty.y {
			empty = pt
		}
	}

	steps := dist(pos, empty)
	mv = move{pod: b.pod, pos: pos, to: empty, steps: steps}
	return mv
}

func getPodMoves(g grid, rooms map[podname]room, pos xy) []xy {
	pts := []xy{}

	if isFinal(g, rooms, pos) {
		return pts
	}

	b := g[pos.y][pos.x]
	if b.hallway {
		return pts
	}

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			pt := xy{x, y}
			b := g[y][x]
			// rooms handled separately
			if b.wall || b.room || b.pod != Empty {
				continue
			}

			if b.hallway && g[y+1][x].room {
				// can't move to spot above room
				continue
			}
			pts = append(pts, pt)
		}
	}

	valid := []xy{}
	for _, pt := range pts {
		if pathClear(g, pos, pt) {
			valid = append(valid, pt)
		}
	}
	return valid
}

func pathClear(g grid, start, end xy) bool {
	trv := start

	for trv.y != hallwayy {
		trv.y--

		if g[trv.y][trv.x].pod != Empty {
			return false
		}
	}

	dir := 1
	if end.x < trv.x {
		dir = -1
	}

	for trv.x != end.x {
		trv.x += dir

		if g[trv.y][trv.x].pod != Empty {
			return false
		}
	}

	dir = 1
	if end.y < trv.y {
		dir = -1
	}

	for trv.y != end.y {
		trv.y += dir
		if g[trv.y][trv.x].pod != Empty {
			return false
		}
	}
	return true
}

func isSolved(g grid, rooms map[podname]room) bool {
	result := true
	for _, r := range rooms {
		for pt := range r.pts {
			b := g[pt.y][pt.x]
			if b.pod.name != r.owner {
				result = false
				break
			}
		}

		if !result {
			break
		}
	}
	return result
}

func parseInput(in input) (grid, map[xy]pod, map[podname]room) {
	g := make(grid, len(in))
	pm := map[xy]pod{}
	roomidx := map[int]podname{3: A, 5: B, 7: C, 9: D}
	plookup := map[podname]int{}
	rooms := map[podname]room{}
	for y, line := range in {
		g[y] = make([]block, len(line))
		for x, c := range line {
			b := block{pod: Empty}
			pt := xy{x, y}
			switch c {
			case '#', ' ':
				b.wall = true
			case '.':
				b.hallway = true

			case 'A', 'B', 'C', 'D':
				if line[x-1] == '#' && line[x+1] == '#' {
					b.room = true
				} else {
					b.hallway = true
				}
				name := podname(string(c))
				energy := 1
				if name == B {
					energy = 10
				} else if name == C {
					energy = 100
				} else if name == D {
					energy = 1000
				}

				idx := plookup[name]

				plookup[name]++

				p := pod{
					name:      name,
					energyuse: energy,
					num:       idx,
				}
				pm[pt] = p
				b.pod = p
			}

			if x > 0 && x < len(line)-1 && !b.wall {
				if line[x-1] == '#' && line[x+1] == '#' {
					b.room = true
					owner, validroom := roomidx[x]
					if _, ok := rooms[owner]; !ok && validroom {
						r := room{
							owner: owner,
							pts:   map[xy]bool{pt: true},
							x:     x,
							y:     y,
						}
						rooms[owner] = r
					} else if validroom {
						r := rooms[owner]
						r.y = y
						r.pts[pt] = true
						rooms[owner] = r
					}

					if !validroom {
						fmt.Println("not valid room encountered", x, y, line[x])
					}
				}
			}

			g[y][x] = b

		}
	}
	return g, pm, rooms
}
