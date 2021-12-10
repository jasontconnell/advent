package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var inputFilename = "input.txt"

type output = int

type state struct {
	moves         int
	floors        [][]string
	elevator      []string
	elevatorLevel int
	score         int
}

type move struct {
	elevatorLevel   int
	pickup, dropoff []string
}

func copyElevator(elevator []string) []string {
	cp := make([]string, len(elevator))
	copy(cp, elevator)
	return cp
}

func copyFloors(floors [][]string) [][]string {
	cp := make([][]string, len(floors))
	for i := 0; i < len(floors); i++ {
		cp[i] = make([]string, len(floors[i]))

		copy(cp[i], floors[i])
	}
	return cp
}

func getKey(floors [][]string, elevator []string, level int) string {
	key := ""

	for floor, r := range floors {
		key += " " + strconv.Itoa(floor) + "| "
		for _, element := range r {
			key += element + " "
		}
	}
	key += "| elevator "
	for _, ev := range elevator {
		key += " " + ev
	}
	key += fmt.Sprintf("| level %d", level)
	return key
}

func isComplete(floors [][]string) bool {
	complete := true
	for i := 0; i < len(floors)-1; i++ {
		complete = complete && len(floors[i]) == 0
	}
	return complete
}

func getScore(floors [][]string) int {
	return len(floors[0]) + 10*len(floors[1]) + 100*len(floors[2]) + 1000*len(floors[3])
}

func main() {
	startTime := time.Now()

	floors := [][]string{}
	floors = append(floors, []string{"PO-G", "TH-G", "TH-M", "PR-G", "RU-G", "RU-M", "CO-G", "CO-M"})
	floors = append(floors, []string{"PO-M", "PR-M"})
	floors = append(floors, []string{})
	floors = append(floors, []string{})

	sort.Strings(floors[0])
	sort.Strings(floors[1])

	p1 := part1(floors)
	p2 := part2(floors)

	fmt.Println("--2016 day 11 solution--")
	fmt.Println("(this one takes a few minutes)")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(floors [][]string) output {
	b, i := simulate(floors)
	if b {
		return i
	}
	return 0
}

func part2(floors [][]string) output {
	floors[0] = append(floors[0], []string{"EL-G", "EL-M", "DI-G", "DI-M"}...)
	sort.Strings(floors[0])
	b, i := simulate(floors)
	if b {
		return i
	}
	return 0
}

func simulate(floors [][]string) (bool, int) {
	floorInit := copyFloors(floors)
	queue := []*state{}
	mvs := getValidMoves(floorInit, []string{}, 0)

	for _, mv := range mvs {
		st := newState(floorInit, []string{}, mv, 1)
		queue = append(queue, st)
	}

	visited := make(map[string]bool)
	moves := 1000
	solved := false

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		key := getKey(cur.floors, cur.elevator, cur.elevatorLevel)
		if _, ok := visited[key]; ok {
			continue
		}
		visited[key] = true

		if isComplete(cur.floors) && cur.moves < moves {
			moves = cur.moves
			solved = true
			continue
		}

		mvs := getValidMoves(cur.floors, cur.elevator, cur.elevatorLevel)

		for _, mv := range mvs {
			st := newState(cur.floors, cur.elevator, mv, cur.moves)

			queue = append(queue, st)
		}
	}

	return solved, moves
}

func newState(floors [][]string, elevator []string, mv move, moves int) *state {
	cp := new(state)
	cp.floors = copyFloors(floors)
	cp.elevator = copyElevator(elevator)
	cp.moves = moves
	cp.elevatorLevel = mv.elevatorLevel

	if mv.pickup != nil && len(mv.pickup) != 0 {
		cp.elevator = addElements(cp.elevator, mv.pickup)
		cp.floors[cp.elevatorLevel] = removeElements(cp.floors[cp.elevatorLevel], mv.pickup)
	}

	if mv.dropoff != nil && len(mv.dropoff) != 0 {
		cp.floors[mv.elevatorLevel] = addElements(cp.floors[mv.elevatorLevel], mv.dropoff)
		cp.elevator = removeElements(cp.elevator, mv.dropoff)
		cp.moves++
	}
	cp.score = getScore(cp.floors)
	return cp
}

func addElements(list []string, elements []string) []string {
	list = append(list, elements...)
	sort.Strings(list)
	return list
}

func removeElements(list []string, elements []string) []string {
	for _, element := range elements {
		for j := len(list) - 1; j >= 0; j-- {
			if list[j] == element {
				list = append(list[:j], list[j+1:]...)
			}
		}
	}
	sort.Strings(list)
	return list
}

func getValidMoves(floors [][]string, elevator []string, currentLevel int) []move {
	mvs := getAllMoves(floors, elevator, currentLevel)
	validMoves := []move{}

	for _, mv := range mvs {
		if mv.dropoff != nil && safeToRemove(elevator, mv.dropoff) && safeToAdd(floors[mv.elevatorLevel], mv.dropoff) {
			validMoves = append(validMoves, mv)
		}
		if mv.pickup != nil && safeToRemove(floors[mv.elevatorLevel], mv.pickup) && safeToAdd(elevator, mv.pickup) {
			validMoves = append(validMoves, mv)
		}
	}

	return validMoves
}

func getAllMoves(floors [][]string, elevator []string, currentLevel int) []move {
	mvs := []move{}

	for _, toFloor := range []int{currentLevel - 1, currentLevel + 1} {
		if toFloor == -1 || toFloor == len(floors) {
			continue
		}
		if toFloor < currentLevel && len(floors[toFloor]) == 0 {
			continue
		}

		// don't drop off 2 below, only move up
		if len(elevator) == 2 && toFloor < currentLevel {
			continue
		}

		if len(elevator) > 0 {
			combos := [][]string{
				{elevator[0]},
			}
			if len(elevator) > 1 && elevator[1] != "" {
				combos = append(combos, [][]string{{elevator[1]}, {elevator[0], elevator[1]}}...)
			}
			for _, elements := range combos {
				mv := move{
					elevatorLevel: toFloor,
					dropoff:       elements,
				}
				mvs = append(mvs, mv)
			}
		}

		if toFloor < currentLevel && len(elevator) < 2 {
			for j := 0; j < len(floors[currentLevel]); j++ {
				elem := floors[currentLevel][j]
				if elem == "" {
					continue
				}
				mv := move{
					elevatorLevel: currentLevel,
					pickup:        []string{elem},
				}
				mvs = append(mvs, mv)
			}
		}

		if toFloor > currentLevel && len(elevator) == 0 {
			for j := 0; j < len(floors[currentLevel]); j++ {
				for x := 0; x < len(floors[currentLevel]); x++ {
					if j == x {
						continue
					}

					left, right := floors[currentLevel][j], floors[currentLevel][x]
					if left == "" || right == "" {
						continue
					}
					pk := []string{left, right}
					sort.Strings(pk)
					mvs = append(mvs, move{pickup: pk, elevatorLevel: currentLevel})
				}
			}
		}
	}
	return mvs
}

func safeToRemove(row []string, elements []string) bool {
	if len(elements) == 0 {
		return true
	}
	cp := make([]string, len(row))
	copy(cp, row)
	for _, element := range elements {
		for j := len(cp) - 1; j >= 0; j-- {
			if cp[j] == element {
				cp = append(cp[:j], cp[j+1:]...)
			}
		}
	}

	return rowValid(cp)
}

func safeToAdd(row []string, elements []string) bool {
	cp := make([]string, len(row))
	copy(cp, row)
	cp = append(cp, elements...)
	sort.Strings(cp)
	return rowValid(cp)
}

func rowValid(row []string) bool {
	result := true
	for _, r := range row {
		if isMicrochip(r) && !containsOwnGenerator(row, r) && containsOtherGenerator(row, r) && !containsDuplicate(row, r) {
			result = false
		}
	}

	return result
}

func isMicrochip(element string) bool {
	return strings.HasSuffix(element, "-M")
}

func containsDuplicate(row []string, element string) bool {
	if len(row) <= 1 {
		return false
	}

	for i := 1; i < len(row); i++ {
		if row[i] == row[i-1] {
			return true
		}
	}
	return false
}

func containsOwnGenerator(row []string, element string) bool {
	result := false
	gen := string(element[0:2]) + "-G"

	for _, r := range row {
		if r == gen {
			result = true
		}
	}

	return result
}

func containsOtherGenerator(row []string, element string) bool {
	result := false
	el := string(element[0:2])

	for _, r := range row {
		if strings.HasSuffix(r, "-G") && !strings.HasPrefix(r, el) {
			result = true
		}
	}
	return result
}
