package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

var inputFilename = "input.txt"

type output = int

var origFloors []map[int]bool

type pqueue []*state

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
	st := x.(*state)
	*pq = append(*pq, st)
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

type state struct {
	moves        int
	floors       []map[int]bool
	elevator     map[int]bool
	currentLevel int
	score        float64
}

type move struct {
	pickup  map[int]bool
	deliver map[int]bool
	level   int
}

func (m move) String() string {
	return fmt.Sprintf("pickup: %v deliver: %v to level %d", sortedKeys(m.pickup), sortedKeys(m.deliver), m.level)
}

func (m move) key(level int) string {
	return fmt.Sprintf("p:%v d:%v l:%d c:%d", sortedKeys(m.pickup), sortedKeys(m.deliver), m.level, level)
}

func copymap(m map[int]bool) map[int]bool {
	cp := make(map[int]bool)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func maplist(list []int) map[int]bool {
	m := map[int]bool{}
	for _, s := range list {
		m[s] = true
	}
	return m
}

func copyFloors(floors []map[int]bool) []map[int]bool {
	cp := make([]map[int]bool, len(floors))
	for i := 0; i < len(floors); i++ {
		cp[i] = copymap(floors[i])
	}
	return cp
}

func getKey(floors []map[int]bool, elevator map[int]bool, level int) string {
	key := ""
	for floor, r := range floors {
		fl := sortedKeys(r)
		key += "|"
		for _, k := range fl {
			s := "G"
			if k < 0 {
				s = "M"
			}
			key += s + " "
		}
		if floor == level {
			el := sortedKeys(elevator)
			key += "e: "
			for _, k := range el {
				s := "G"
				if k < 0 {
					s = "M"
				}
				key += s + " "
			}
		}
	}
	return key
}

func sortedKeys(m map[int]bool) []int {
	list := []int{}
	for k := range m {
		list = append(list, k)
	}
	sort.Ints(list)
	return list
}

func isComplete(floors []map[int]bool, elevator map[int]bool, level int) bool {
	if level != len(floors)-1 {
		return false
	}
	if len(elevator) > 0 {
		return false
	}
	complete := true
	for i, m := range floors {
		if i != len(floors)-1 && len(m) > 0 {
			complete = false
			break
		}
	}
	return complete
}

func getScore(moves int, floors []map[int]bool, level int) float64 {
	dist := len(floors[0])*1000 + len(floors[1])*100 + len(floors[2])*10 + len(floors[3])
	return float64(dist)
}

func main() {
	startTime := time.Now()

	origFloors = make([]map[int]bool, 4)
	// origFloors[0] = map[int]bool{"PO-G": true, "TH-G": true, "TH-M": true, "PR-G": true, "RU-G": true, "RU-M": true, "CO-G": true, "CO-M": true}
	// origFloors[1] = map[int]bool{"PO-M": true, "PR-M": true}
	origFloors[0] = map[int]bool{1: true, 2: true, -2: true, 3: true, 4: true, -4: true, 5: true, -5: true}
	origFloors[2] = map[int]bool{-1: true, -3: true}
	origFloors[3] = map[int]bool{}

	runExample := len(os.Args) > 1 && os.Args[1] == "example"
	if runExample {
		fmt.Println("using example input")
		example := []map[int]bool{{-1: true, -2: true}, {1: true}, {2: true}}
		origFloors[0] = example[0]
		origFloors[1] = example[1]
		origFloors[2] = example[2]
	}

	p1 := part1(origFloors)
	p2 := part2(origFloors)

	fmt.Println("--2016 day 11 solution--")
	fmt.Println("(this one takes a few minutes)")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(floors []map[int]bool) output {
	b, i := simulate(floors)
	if b {
		return i
	}
	return 0
}

func part2(floors []map[int]bool) output {
	for _, s := range []int{6, -6, 7, -7} {
		floors[0][s] = true
	}
	b, i := simulate(floors)
	if b {
		return i
	}
	return 0
}

func simulate(floors []map[int]bool) (bool, int) {
	floorInit := copyFloors(floors)
	queue := pqueue{}

	initial := newState(floorInit, map[int]bool{}, 0)
	heap.Init(&queue)
	queue.Push(initial)

	visited := make(map[string]int)
	moves := math.MaxInt32
	solved := false

	for queue.Len() > 0 {
		cur := (heap.Pop(&queue)).(*state)

		if isComplete(cur.floors, cur.elevator, cur.currentLevel) && cur.moves < moves {
			moves = cur.moves
			solved = true
			continue
		}

		if cur.moves > moves {
			continue
		}

		key := getKey(cur.floors, cur.elevator, cur.currentLevel)
		if _, ok := visited[key]; ok {
			continue
		}
		visited[key] = cur.moves

		mvs := getValidMoves(cur.floors, cur.elevator, cur.currentLevel)
		for _, mv := range mvs {
			transit := transitState(cur, mv)
			queue = append(queue, transit)
		}
	}

	return solved, moves
}

func copyState(st *state) *state {
	cp := new(state)
	cp.floors = copyFloors(st.floors)
	cp.elevator = copymap(st.elevator)
	cp.currentLevel = st.currentLevel
	cp.moves = st.moves
	return cp
}

func transitState(st *state, mv move) *state {
	cp := copyState(st)

	cp.floors[cp.currentLevel], cp.elevator = transfer(cp.floors[cp.currentLevel], cp.elevator, mv.pickup)
	cp.elevator, cp.floors[mv.level] = transfer(cp.elevator, cp.floors[mv.level], mv.deliver)

	cp.moves++
	cp.currentLevel = mv.level
	cp.score = getScore(cp.moves, cp.floors, cp.currentLevel)
	return cp
}

func newState(floors []map[int]bool, elevator map[int]bool, level int) *state {
	cp := new(state)
	cp.moves = 0
	cp.floors = copyFloors(floors)
	cp.elevator = copymap(elevator)
	cp.currentLevel = level
	return cp
}

func transfer(from map[int]bool, to map[int]bool, elements map[int]bool) (map[int]bool, map[int]bool) {
	for element := range elements {
		if _, ok := from[element]; ok {
			if _, ok := to[element]; !ok {
				delete(from, element)
				to[element] = true
			}
		}
	}
	return from, to
}

func getValidMoves(floors []map[int]bool, elevator map[int]bool, currentLevel int) []move {
	mvs := getAllMoves(floors, elevator, currentLevel)
	validMoves := []move{}

	for _, mv := range mvs {
		valid := true

		if !safeToRemove(floors[currentLevel], mv.pickup) || !safeToAdd(floors[mv.level], mv.deliver) {
			valid = false
		}

		if valid {
			validMoves = append(validMoves, mv)
		}
	}
	return validMoves
}

func getAllMoves(floors []map[int]bool, elevator map[int]bool, currentLevel int) []move {
	elevatorKeys := sortedKeys(elevator)
	floorKeys := sortedKeys(floors[currentLevel])

	mvs := []move{}
	for _, level := range []int{currentLevel + 1, currentLevel - 1} {
		if level < 0 || level == len(floors) {
			continue
		}
		if level < currentLevel && len(elevatorKeys) == 2 {
			continue
		}

		toElevator := [][]int{}
		if len(floorKeys) > 0 {
			m := map[int]bool{}
			for j := 0; j < len(floorKeys); j++ {
				left := floorKeys[j]
				for r := j + 1; r < len(floorKeys); r++ {
					right := floorKeys[r]

					if len(elevatorKeys) == 0 && level > currentLevel {
						toElevator = append(toElevator, []int{left, right})
					}
					if _, ok := m[right]; !ok && len(elevatorKeys) < 2 {
						toElevator = append(toElevator, []int{right})
						m[right] = true
					}
				}
				if _, ok := m[left]; !ok && len(elevatorKeys) < 2 {
					toElevator = append(toElevator, []int{left})
					m[left] = true
				}
			}
		}

		for i := 0; i < len(toElevator); i++ {
			em := toElevator[i]

			if len(em) > 0 {
				left := em[0]
				if len(em) == 2 {
					right := em[1]

					mvs = append(mvs,
						move{
							level:   level,
							pickup:  map[int]bool{left: true, right: true},
							deliver: map[int]bool{left: true, right: true},
						})

					mvs = append(mvs,
						move{
							level:   level,
							pickup:  map[int]bool{left: true, right: true},
							deliver: map[int]bool{left: true},
						})
					mvs = append(mvs,
						move{
							level:   level,
							pickup:  map[int]bool{left: true, right: true},
							deliver: map[int]bool{right: true},
						})
					mvs = append(mvs,
						move{
							level:   level,
							pickup:  map[int]bool{right: true},
							deliver: map[int]bool{right: true},
						})
				}

				mvs = append(mvs,
					move{
						level:   level,
						pickup:  map[int]bool{left: true},
						deliver: map[int]bool{left: true},
					})
			}
		}
	}

	return mvs
}

func safeToRemove(row map[int]bool, elements map[int]bool) bool {
	if len(elements) == 0 {
		return true
	}
	cp := make(map[int]bool, len(row))
	for element := range row {
		cp[element] = true
	}
	for element := range elements {
		delete(cp, element)
	}
	return rowValid(cp)
}

func safeToAdd(row map[int]bool, elements map[int]bool) bool {
	if len(elements) == 0 {
		return true
	}
	cp := make(map[int]bool, len(row))
	for element := range row {
		cp[element] = true
	}
	for element := range elements {
		cp[element] = true
	}
	return rowValid(cp)
}

func rowValid(row map[int]bool) bool {
	result := true
	for element := range row {
		if isMicrochip(element) && !containsOwnGenerator(row, element) && containsOtherGenerator(row, element) {
			result = false
		}
	}
	return result
}

func isMicrochip(element int) bool {
	return element < 0
}

func containsOwnGenerator(row map[int]bool, element int) bool {
	_, ok := row[-element]
	return ok
}

func containsOtherGenerator(row map[int]bool, element int) bool {
	result := false

	for r := range row {
		if r > 0 && r != -element {
			result = true
		}
	}
	return result
}
