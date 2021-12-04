package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "04.txt"
var size int = 5

type xy struct {
	x, y int
}

type square struct {
	point xy
	num   int
}

type board struct {
	entries map[xy]square
	valmap  map[int]square
}

func (b board) String() string {
	s := ""
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			e, _ := b.entries[xy{x, y}]
			s += fmt.Sprintf("%v ", e.num)
		}
		s += "\n"
	}

	return s
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

	nums, boards := getInput(lines)

	p1 := part1(nums, boards)
	p2 := part2(nums, boards)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(nums []int, boards []board) int {
	val := 0
	done := false
	for i := 0; i < len(nums) && !done; i++ {
		for _, b := range boards {
			if checkBoard(b, nums, i) {
				sum := getUnmarked(b, nums, i)
				val = sum * nums[i-1]
				done = true
				break
			}
		}
	}
	return val
}

func part2(nums []int, boards []board) int {
	winners := make(map[int]board)
	last := 0
	for i := 0; i < len(nums); i++ {
		for ib, b := range boards {
			if _, ok := winners[ib]; ok {
				continue
			}

			if checkBoard(b, nums, i) {
				winners[ib] = b
				sum := getUnmarked(b, nums, i)
				last = sum * nums[i-1]
			}
		}
	}
	return last
}

func getUnmarked(b board, nums []int, itr int) int {
	m := make(map[int]int)
	for _, n := range nums[:itr] {
		m[n] = n
	}

	val := 0
	for n, _ := range b.valmap {
		if _, ok := m[n]; !ok {
			val += n
		}
	}
	return val
}

func checkBoard(b board, nums []int, itr int) bool {
	matches := make(map[xy]bool)
	for _, n := range nums[:itr] {
		if sq, ok := b.valmap[n]; ok {
			matches[sq.point] = true
		}
	}

	match := false
	for p, _ := range matches {
		rm := checkRow(matches, p)

		if rm {
			match = true
			break
		}

		cm := checkCol(matches, p)
		if cm {
			match = true
			break
		}
	}

	return match
}

func checkRow(matches map[xy]bool, p xy) bool {
	val := 0
	for i := 0; i < size; i++ {
		cp := xy{i, p.y}
		if _, ok := matches[cp]; ok {
			val++
		}
	}
	return val == size
}

func checkCol(matches map[xy]bool, p xy) bool {
	val := 0
	for i := 0; i < size; i++ {
		cp := xy{p.x, i}
		if _, ok := matches[cp]; ok {
			val++
		}
	}
	return val == size
}

func getInput(lines []string) ([]int, []board) {
	nums := getNums(lines[0])

	boards := []board{}

	var cur board
	var row int
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			cur = board{}
			cur.entries = make(map[xy]square)
			cur.valmap = make(map[int]square)
			row = 0
			boards = append(boards, cur)
			continue
		}

		squares := getSquares(lines[i], row)
		for _, sq := range squares {
			cur.valmap[sq.num] = sq
			cur.entries[sq.point] = sq
		}
		row++
	}

	return nums, boards
}

func getNums(line string) []int {
	vals := strings.Split(line, ",")
	ret := []int{}
	for _, v := range vals {
		i, _ := strconv.Atoi(v)
		ret = append(ret, i)
	}
	return ret
}

func getSquares(line string, row int) []square {
	trimmed := strings.Trim(line, " ")
	split := strings.Fields(trimmed)

	squares := []square{}
	for x, sp := range split {
		v, _ := strconv.Atoi(sp)
		sq := square{num: v, point: xy{x, row}}
		squares = append(squares, sq)
	}
	return squares
}
