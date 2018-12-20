package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

var input = "15.txt"

const (
	defaultHP  = 200
	defaultAtk = 3
)

const (
	Open int = iota
	Wall
)

const (
	Elf int = iota
	Goblin
)

type xy struct {
	x, y int
}

func (v xy) String() string {
	return fmt.Sprintf("(%d, %d)", v.x, v.y)
}

type path struct {
	xy
	block int // open or blocked
	u     unit
}

type unit struct {
	xy
	id   int
	race int
	atk  int
	hp   int
	dead bool
}

func (u unit) String() string {
	race := "Elf"
	if u.race == Goblin {
		race = "Goblin"
	}
	return fmt.Sprintf("%v - %d - %s  HP: %d", u.xy, u.id, race, u.hp)
}

type state struct {
	xy
	moves []xy
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	grid := [][]path{}
	units := []unit{}
	y := 0
	for scanner.Scan() {
		var txt = scanner.Text()
		p, u := getLine(y, txt, defaultHP, defaultAtk)
		grid = append(grid, p)
		units = append(units, u...)
		y++
	}

	for i := 0; i < len(units); i++ {
		units[i].id = i
	}

	// print(units, grid)

	units, round := sim(units, grid)

	sum := 0
	for _, u := range units {
		sum += u.hp
	}
	fmt.Println("Part 1:", sum*round)
	fmt.Println("Time", time.Since(startTime))
}

func printUnits(units []unit) {
	s := ""
	for _, u := range units {
		s += fmt.Sprintf("[%v] :: ", u)
	}
	fmt.Println(s)
}

func print(units []unit, grid [][]path) {
	umap := unitMap(units)

	for y := 0; y < len(grid); y++ {
		line := ""
		for x := 0; x < len(grid[y]); x++ {
			key := xy{x: x, y: y}
			g := grid[y][x]
			ch := '.'
			if g.block == Wall {
				ch = '#'
			}
			if u, ok := umap[key]; ok {
				ch = 'G'
				if u.race == Elf {
					ch = 'E'
				}
				//ch = rune(48 + u.id)
				if u.dead {
					ch = '.'
				}
			}
			line += string(ch)
		}
		fmt.Println(line)
	}
	printUnits(units)
	fmt.Println("--------------------------------------------")
}

func sim(units []unit, grid [][]path) ([]unit, int) {
	done := false
	roundNum := 0
	for !done {
		units = turn(units, grid)
		units = bringOutYourDead(units)
		enlist := enemies(units[0], units)
		done = len(enlist) == 0
		//print(units, grid)
		if !done {
			roundNum++
		}
	}

	return units, roundNum
}

func bringOutYourDead(units []unit) []unit {
	for i := len(units) - 1; i >= 0; i-- {
		if units[i].dead {
			units = append(units[:i], units[i+1:]...)
		}
	}
	return units
}

func attack(u unit, units []unit) (unit, int) { // returns unit and index
	atk, index := unit{id: -1, hp: 201}, -1
	enmap := unitMap(enemies(u, units))

	sur := surrounding(u.xy)
	for _, s := range sur {
		if en, ok := enmap[s]; ok && en.hp < atk.hp {
			atk = en
		}
	}

	if atk.id == -1 {
		return unit{}, -1
	}

	for i, un := range units {
		if un.id == atk.id {
			index = i
			break
		}
	}

	if index != -1 {
		atk.hp = atk.hp - u.atk
		if atk.hp < 0 {
			atk.dead = true
		}
	}
	return atk, index
}

func turn(units []unit, grid [][]path) []unit {
	usort := sortUnits(units)

	for i := 0; i < len(usort); i++ {
		u := usort[i]
		atk := canAttack(u, units)

		if !atk {
			n := getNext(u, usort, grid)
			if n.x != -1 && n.y != -1 {
				u = move(u, n, units, grid)
			}
		}

		dmgd, index := attack(u, usort)
		if index != -1 {
			usort[index] = dmgd
		}

		usort[i] = u
		// print(usort, grid)
	}
	return usort
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func move(u unit, to xy, units []unit, grid [][]path) unit {
	u.x = to.x
	u.y = to.y

	return u
}

func sortUnits(units []unit) []unit {
	s := func(i, j int) bool {
		dy := units[i].y - units[j].y
		l := dy < 0

		if dy == 0 {
			l = units[i].x < units[j].x
		}
		return l
	}

	sort.Slice(units, s)

	return units
}

func distance(p1, p2 xy) int {
	dx := abs(p1.x - p2.x)
	dy := abs(p1.y - p2.y)
	return dx + dy
}

func getNext(u unit, units []unit, grid [][]path) xy {
	enlist := enemies(u, units)
	goals := []xy{}
	for _, e := range enlist {
		spots := availableSpots(e, units, grid)
		for _, s := range spots {
			goals = append(goals, s)
		}
	}

	if len(goals) == 0 {
		return u.xy
	}

	open := getOpen(units, grid)

	min := 10000
	mv := xy{-1, -1}
	for _, g := range goals {
		visited := make(map[xy]bool)
		shortest := getPath(u.xy, g, open, visited)
		if len(shortest.moves) > 0 {
			final := shortest.moves[len(shortest.moves)-1]
			dist := distance(final, u.xy)
			if dist < min {
				min = dist
				mv = shortest.moves[0]
			}
		}
	}
	return mv
}

func enemies(u unit, units []unit) []unit {
	en := []unit{}
	for _, uu := range units {
		if uu.race != u.race && !uu.dead {
			en = append(en, uu)
		}
	}
	return en
}

func getPath(from, to xy, open []xy, visited map[xy]bool) state {
	queue := []state{}
	queue = append(queue, state{xy: from, moves: []xy{}})
	solves := []state{}
	minsolve := 10000

	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		mvs := getMoves(s.xy, open)

		for _, mv := range mvs {
			mvstate := state{moves: append(s.moves, mv), xy: mv}
			if mv.x == to.x && mv.y == to.y {
				if len(mvstate.moves) < minsolve {
					minsolve = len(mvstate.moves)
					solves = append([]state{mvstate}, solves...) // move shortest path to front
				}
			}

			if _, ok := visited[mv]; !ok {
				visited[mv] = true

				if len(mvstate.moves)+1 < minsolve {
					queue = append(queue, mvstate)
				}
			}
		}
	}

	mv := state{}
	if len(solves) > 0 {
		mv = solves[0]
	}
	return mv
}

func getMoves(from xy, open []xy) []xy {
	xym := xyMap(open)
	moves := []xy{}

	for _, p := range surrounding(from) {
		if _, ok := xym[p]; ok {
			moves = append(moves, p)
		}
	}
	return moves
}

func getOpen(units []unit, grid [][]path) []xy {
	umap := unitMap(units)
	open := []xy{}
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			key := xy{x: x, y: y}
			if _, ok := umap[key]; !ok && grid[y][x].block == Open {
				open = append(open, key)
			}
		}
	}
	return open
}

func availableSpots(u unit, units []unit, grid [][]path) []xy {
	avail := []xy{}
	umap := unitMap(units)

	for _, point := range surrounding(u.xy) {
		g := grid[point.y][point.x]
		if _, ok := umap[point]; !ok && g.block == Open {
			avail = append(avail, point)
		}
	}
	return avail
}

func canAttack(u unit, units []unit) bool {
	emap := unitMap(enemies(u, units))
	for _, point := range surrounding(u.xy) {
		if _, ok := emap[point]; ok {
			return true
		}
	}
	return false
}

func surrounding(p xy) []xy {
	return []xy{
		xy{x: p.x, y: p.y - 1},
		xy{x: p.x - 1, y: p.y},
		xy{x: p.x + 1, y: p.y},
		xy{x: p.x, y: p.y + 1},
	}
}

func unitMap(units []unit) map[xy]unit {
	umap := make(map[xy]unit)
	for _, u := range units {
		if u.dead {
			continue
		}
		umap[u.xy] = u
	}
	return umap
}

func xyMap(xys []xy) map[xy]bool {
	xym := make(map[xy]bool)
	for _, p := range xys {
		xym[p] = true
	}
	return xym
}

func getLine(y int, txt string, hp, atk int) ([]path, []unit) {
	units := []unit{}
	paths := []path{}
	for x, c := range txt {
		coords := xy{x: x, y: y}
		p := path{xy: coords}
		switch c {
		case '#':
			p.block = Wall
		case '.':
			p.block = Open
		case 'E', 'G':
			p.block = Open
			r := Elf
			if c == 'G' {
				r = Goblin
			}
			u := unit{xy: coords, race: r, hp: hp, atk: atk}
			units = append(units, u)
		}
		paths = append(paths, p)
	}

	return paths, units
}
