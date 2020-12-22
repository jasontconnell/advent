package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "20.txt"

type edge int

const (
	top edge = iota
	bottom
	left
	right
	undefined
)

var edges []edge = []edge{top, bottom, left, right}

func (e edge) opposite() edge {
	switch e {
	case top:
		return bottom
	case bottom:
		return top
	case left:
		return right
	case right:
		return left
	}
	return undefined
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
	d := 10

	ms := readGrids(lines, d)
	p1 := solve(ms, d)

	fmt.Println("Part 1:", p1)

	fmt.Println("Time", time.Since(startTime))
}

func solve(m map[int][][]bool, d int) int {
	blocks := make(map[int]map[edge]int)

	unmatched := []int{}

	for k := range m {
		unmatched = append(unmatched, k)
		blocks[k] = make(map[edge]int)
	}

	for len(unmatched) > 0 {
		k := unmatched[0]
		unmatched = unmatched[1:]
		exclude := []int{k}

		for _, found := range blocks[k] {
			exclude = append(exclude, found)
		}

		mid, side := getMatch(m[k], m, d, exclude)

		if mid != 0 && side != undefined {
			blocks[k][side] = mid
			if len(blocks[k]) < 4 {
				unmatched = append(unmatched, k)
			}
		}
	}
	prod := 1
	for k, v := range blocks {
		if len(v) == 2 {
			prod *= k
		}
	}
	return prod
}

// rearrange the image after all things are lined up
func rearrange(m map[int][][]bool, d int) {

}

func getMatch(list [][]bool, m map[int][][]bool, d int, exclude []int) (int, edge) {
	matched := 0
	exmap := make(map[int]int)
	for _, ex := range exclude {
		exmap[ex] = ex
	}
	var side edge
	for k := range m {
		if _, ok := exmap[k]; ok {
			continue
		}
		test := m[k]

		ok, s := isMatch(list, test, d)

		if ok {
			//m[k] = test // new orientation?
			matched = k
			side = s
			break
		}
	}
	return matched, side
}

func isMatch(list [][]bool, test [][]bool, d int) (bool, edge) {
	b := false
	e := top
	i := 0
	for i < 4 {
		b, e = checkEdges(list, test, d)

		if !b {
			test = flipVertical(test, d)
			b, e = checkEdges(list, test, d)

			if !b {
				test = flipVertical(test, d)
				test = rotate(test, d)
			}
		}
		i++
	}
	return b, e
}

func checkEdges(list [][]bool, test [][]bool, d int) (bool, edge) {
	edgemap := make(map[edge]int)

	for i := 0; i < d; i++ {
		if list[i][0] == test[i][0] {
			edgemap[left]++
		}

		if list[0][i] == test[0][i] {
			edgemap[top]++
		}

		if list[d-1][i] == test[d-1][i] {
			edgemap[bottom]++
		}

		if list[i][d-1] == test[i][d-1] {
			edgemap[right]++
		}
	}

	matched := false
	side := undefined
	for k, v := range edgemap {
		if v == d {
			matched = true
			side = k
			break
		}
	}
	return matched, side
}

func printGrid(m [][]bool, d int) {
	for y := 0; y < d; y++ {
		for x := 0; x < d; x++ {
			on := m[y][x]
			if on {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func flipVertical(m [][]bool, d int) [][]bool {
	m2 := make([][]bool, d)
	for y := 0; y < d; y++ {
		m2[y] = make([]bool, d)
	}
	for y := 0; y < d; y++ {
		for x := 0; x < d; x++ {
			m2[d-y-1][x] = m[y][x]
		}
	}
	return m2
}

func rotate(m [][]bool, d int) [][]bool {
	m2 := make([][]bool, d)
	for y := 0; y < d; y++ {
		m2[y] = make([]bool, d)
		for x := 0; x < d; x++ {
			m2[y][x] = m[d-x-1][y]
		}
	}
	return m2
}

func readGrids(lines []string, d int) map[int][][]bool {
	m := make(map[int][][]bool)
	curId := 0
	curm := make([][]bool, d)
	x, y := 0, 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "Tile ") {
			curId, _ = strconv.Atoi(line[5 : len(line)-1])
			curm = make([][]bool, d)
			m[curId] = curm
			y = 0
			continue
		}

		x = 0
		curm[y] = make([]bool, d)
		for _, ch := range line {
			curm[y][x] = ch == '#'
			x++
		}
		y++
	}
	return m
}
