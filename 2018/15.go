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
	dist  int
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

	p1 := sim(units, grid, 3, false)
	var p2 int
	atk := 4
	for p2 == 0 {
		p2 = sim(units, grid, atk, true)
		atk++
	}

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func sim(orig []unit, grid [][]path, atk int, endOnElfDeath bool) int {
	cloned := append([]unit{}, orig...)

	for i := 0; i < len(cloned); i++ {
		if cloned[i].race == Elf {
			cloned[i].atk = atk
		}
	}

	done := false
	elfDied := false
	roundNum := 0
	var fullRound bool
	for !done {
		cloned, fullRound, elfDied = turn(cloned, grid)
		if endOnElfDeath && elfDied {
			return 0
		}
		cloned = bringOutYourDead(cloned)
		enlist := enemies(cloned[0], cloned)
		if fullRound {
			roundNum++
		}
		done = len(enlist) == 0
	}

	sum := 0
	for _, u := range cloned {
		sum += u.hp
	}

	return roundNum * sum
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
	enlist := enemies(u, units)
	enmap := unitMap(enlist)

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
		if atk.hp <= 0 {
			atk.dead = true
		}
	}
	return atk, index
}

func turn(units []unit, grid [][]path) ([]unit, bool, bool) { // whether a full round completed and whether an elf died
	fullRound := true
	elfDied := false
	usort := sortUnits(units)

	for i := 0; i < len(usort); i++ {
		u := usort[i]
		if u.dead {
			continue
		}

		enlist := enemies(u, usort)
		if len(enlist) == 0 && i < len(usort)-1 {
			fullRound = false
		}

		atk := canAttack(u, enlist)
		if !atk {
			n, m := getNext(u, usort, grid)
			if m {
				u = move(u, n)
			}
			atk = canAttack(u, enlist)
		}

		if atk {
			dmgd, index := attack(u, usort)
			usort[index] = dmgd

			if !elfDied {
				elfDied = dmgd.race == Elf && dmgd.dead
			}
		}

		usort[i] = u
	}
	return usort, fullRound, elfDied
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func move(u unit, to xy) unit {
	u.x = to.x
	u.y = to.y

	return u
}

func sortXY(xys []xy) []xy {
	s := func(i, j int) bool {
		dy := xys[i].y - xys[j].y
		l := dy < 0

		if dy == 0 {
			l = xys[i].x < xys[j].x
		}
		return l
	}

	sort.Slice(xys, s)

	return xys
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

func getNext(u unit, units []unit, grid [][]path) (xy, bool) {
	enlist := enemies(u, units)
	goals := []xy{}
	open := getOpen(units, grid)

	for _, e := range enlist {
		best := bestSpots(u.xy, e.xy, open)
		goals = append(goals, best...)
	}

	if len(goals) == 0 {
		return xy{}, false
	}

	mv := xy{}
	doMove := false
	min := 1000
	for _, g := range goals {
		visited := make(map[xy]bool)
		shortest := getPath(u.xy, g, open, visited)
		if len(shortest.moves) > 0 {
			doMove = true
			dist := len(shortest.moves)
			if dist < min {
				min = dist
				mv = shortest.moves[0]
			} else if dist == min {
				sorted := sortXY([]xy{mv, shortest.moves[0]})
				mv = sorted[0]
			}
		}
	}

	return mv, doMove
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
	queue = append(queue, state{xy: from, moves: []xy{}, dist: distance(from, to)})
	solves := []state{}
	minsolve := 10000
	omap := xyMap(open)

	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		mvs := getMoves(s.xy, omap)

		for _, mv := range mvs {
			dist := distance(mv, to)
			mvstate := state{moves: append(s.moves, mv), xy: mv, dist: dist}
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

func getMoves(from xy, omap map[xy]bool) []xy {
	moves := []xy{}
	for _, p := range surrounding(from) {
		if _, ok := omap[p]; ok {
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

func bestSpots(from, to xy, open []xy) []xy {
	omap := xyMap(open)
	mindist := 1000
	dmap := make(map[int][]xy)

	for _, point := range surrounding(to) {
		if _, isOpen := omap[point]; isOpen {
			dist := distance(from, to)
			if dist < mindist {
				mindist = dist
			}
			dmap[dist] = append(dmap[dist], point)
		}
	}
	if list, ok := dmap[mindist]; ok && len(list) > 0 {
		return list
	}
	return []xy{}
}

func canAttack(u unit, enlist []unit) bool {
	emap := unitMap(enlist)
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
