package main

import (
	"fmt"
	"time"

	//"regexp"
	"sort"
	"strconv"
	"strings"
	//"math"
	//"math/rand"
)

type Move struct {
	Elevator   []string
	StartFloor int
	EndFloor   int
}

type Layout struct {
	PreviousState   *Layout
	ElevatorFloor   int
	Score           int
	Floors          [][]string
	MaxLevel        int
	ModifiedLayouts LayoutList
}

type LayoutList []*Layout

func (p LayoutList) Len() int { return len(p) }
func (p LayoutList) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}
func (p LayoutList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SortLayouts(list LayoutList) LayoutList {
	sort.Sort(sort.Reverse(list))

	return list
}

func NewLayout(elevatorFloor, maxLevel int) *Layout {
	c := new(Layout)
	c.MaxLevel = maxLevel
	c.ElevatorFloor = elevatorFloor
	c.PreviousState = nil
	return c
}

func (layout *Layout) Copy() *Layout {
	cp := NewLayout(layout.ElevatorFloor, layout.MaxLevel)
	for i := 0; i < len(layout.Floors); i++ {
		cp.Floors = append(cp.Floors, make([]string, len(layout.Floors[i])))
		copy(cp.Floors[i], layout.Floors[i])
	}
	cp.Score = layout.Score
	cp.PreviousState = nil
	return cp
}

func (layout *Layout) GetKey() string {
	key := ""

	for floor, r := range layout.Floors {
		key += " " + strconv.Itoa(floor) + "| "
		for _, element := range r {
			key += element + " "
		}
	}
	return key
}

func (layout Layout) GetScore() int {
	return len(layout.Floors[0]) + len(layout.Floors[1])*10 + len(layout.Floors[2])*100 + len(layout.Floors[3])*1000
}

func main() {
	startTime := time.Now()

	layout := NewLayout(0, 4)
	layout.Floors = append(layout.Floors, []string{"PO-G", "TH-G", "TH-M", "PR-G", "RU-G", "RU-M", "CO-G", "CO-M", "EL-G", "EL-M", "DI-G", "DI-M"})
	layout.Floors = append(layout.Floors, []string{"PO-M", "PR-M"})
	layout.Floors = append(layout.Floors, []string{})
	layout.Floors = append(layout.Floors, []string{})

	moveTotals := []int{}
	attempts := make(map[string]int)

	solved, moves := simulate(layout, 0, 0, layout.MaxLevel, attempts)

	if solved {
		moveTotals = append(moveTotals, moves)
		fmt.Println("got a solution", moves)
	}

	sort.Ints(moveTotals)
	fmt.Println(moveTotals)
	fmt.Println(attempts)

	fmt.Println("Time", time.Since(startTime))
}

func copyAttempts(attempts map[string]int) map[string]int {
	cp := make(map[string]int)
	for k, v := range attempts {
		cp[k] = v
	}
	return cp
}

func simulate(initialLayout *Layout, currentLevel, moveCount, maxLevel int, attempts map[string]int) (bool, int) {
	solutionFound := false
	localCount := moveCount
	layout := initialLayout.Copy()
	moves := getValidMoves(layout, currentLevel, maxLevel)

	for _, mv := range moves {
		movedLayout := move(layout, maxLevel, mv.Elevator, mv.StartFloor, mv.EndFloor)
		key := movedLayout.GetKey()

		solutionFound = layout.Score == 14000
		if _, exists := attempts[key]; !exists && !solutionFound {
			attempts[key] = layout.Score
			layout.ModifiedLayouts = append(layout.ModifiedLayouts, movedLayout)
		}

		if solutionFound {
			fmt.Println("SOLUTION FOUND IN", localCount, "MOVES")
		}
	}

	if solutionFound {
		return solutionFound, localCount
	}

	SortLayouts(layout.ModifiedLayouts)

	for _, modifiedLayout := range layout.ModifiedLayouts {
		attemptsCopy := copyAttempts(attempts) // prevent duplicates down each path
		solutionFound, localCount = simulate(modifiedLayout, modifiedLayout.ElevatorFloor, moveCount+1, maxLevel, attemptsCopy)
	}

	return solutionFound, localCount
}

func getValidMoves(layout *Layout, currentLevel, maxLevel int) []Move {
	moves := getAllMoves(layout, currentLevel, maxLevel)
	validMoves := []Move{}

	for _, move := range moves {
		if safeToRemove(layout.Floors[move.StartFloor], move.Elevator) && safeToAdd(layout.Floors[move.EndFloor], move.Elevator) {
			validMoves = append(validMoves, move)
		}
	}

	return validMoves
}

func getAllMoves(layout *Layout, currentLevel, maxLevel int) []Move {
	moves := []Move{}

	for _, toFloor := range []int{currentLevel - 1, currentLevel + 1} {
		if toFloor == -1 || toFloor == maxLevel {
			continue
		}
		if toFloor == 0 && len(layout.Floors[toFloor]) == 0 {
			continue
		}
		if toFloor == 1 && len(layout.Floors[0]) == 0 && len(layout.Floors[1]) == 0 {
			continue
		}

		if toFloor < currentLevel {
			for j := 0; j < len(layout.Floors[currentLevel]); j++ {
				moves = append(moves, Move{Elevator: []string{layout.Floors[currentLevel][j]}, StartFloor: currentLevel, EndFloor: toFloor})
			}
		}

		if toFloor > currentLevel {
			for j := 0; j < len(layout.Floors[currentLevel]); j++ {
				for x := 0; x < len(layout.Floors[currentLevel]); x++ {
					if j == x {
						continue
					}

					left, right := layout.Floors[currentLevel][j], layout.Floors[currentLevel][x]
					moves = append(moves, Move{Elevator: []string{left, right}, StartFloor: currentLevel, EndFloor: toFloor})
				}
			}
		}
	}
	return moves
}

func move(initialLayout *Layout, maxLevel int, elevator []string, startLevel, endLevel int) *Layout {
	layout := initialLayout.Copy()

	for i := 0; i < len(elevator); i++ {
		for j := 0; j < len(layout.Floors[startLevel]); j++ {
			if layout.Floors[startLevel][j] == elevator[i] {
				layout.Floors[startLevel] = append(layout.Floors[startLevel][:j], layout.Floors[startLevel][j+1:]...)
			}
		}

		layout.Floors[endLevel] = append(layout.Floors[endLevel], elevator[i])
	}

	layout.ElevatorFloor = endLevel
	layout.Score = layout.GetScore()

	sort.Strings(layout.Floors[endLevel])
	sort.Strings(layout.Floors[startLevel])

	layout.PreviousState = initialLayout

	return layout
}

func safeToRemove(row []string, elevator []string) bool {
	cp := make([]string, len(row))
	copy(cp, row)
	for _, element := range elevator {
		for j := len(cp) - 1; j >= 0; j-- {
			if cp[j] == element {
				cp = append(cp[:j], cp[j+1:]...)
			}
		}
	}

	return rowValid(cp)
}

func safeToAdd(row []string, elevator []string) bool {
	cp := make([]string, len(row))
	copy(cp, row)
	cp = append(cp, elevator...)

	return rowValid(cp)
}

func rowValid(row []string) bool {
	result := true
	for _, r := range row {
		if isMicrochip(r) && !containsOwnGenerator(row, r) && containsOtherGenerator(row, r) {
			result = false
		}
	}

	return result
}

func isMicrochip(element string) bool {
	return strings.HasSuffix(element, "-M")
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
