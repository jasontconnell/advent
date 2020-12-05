package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

var input = "05.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	ops := [][]int{}

	for scanner.Scan() {
		var txt = scanner.Text()
		ops = append(ops, getOps(txt))
	}

	allseats := make(map[int]string)

	for i := 1; i < 127; i++ {
		for j := 1; j < 7; j++ {
			s := strconv.Itoa(i+1) + "x" + strconv.Itoa(j+1)
			allseats[(i*8)+j] = s
		}
	}

	max := 0
	for _, ticket := range ops {
		row, col := findSeat(ticket, 128, 8)

		id := (row * 8) + col
		max = int(math.Max(float64(max), float64(id)))

		//s := strconv.Itoa(row) + "x" + strconv.Itoa(col)
		delete(allseats, id)
	}
	fmt.Println("Part 1, max id:", max)

	list := []int{}
	for k, _ := range allseats {
		list = append(list, k)
	}

	sort.Ints(list)
	for _, id := range list {
		x1, x2 := id-1, id+1
		_, x1ok := allseats[x1]
		_, x2ok := allseats[x2]

		if !x1ok && !x2ok {
			fmt.Println("Part 2, ticket id:", id)
			break
		}
	}
	fmt.Println("Time", time.Since(startTime))
}

func findSeat(ops []int, rows, cols int) (int, int) {
	minr, maxr := 0, rows
	minc, maxc := 0, cols
	for i, x := range ops {
		// rows
		if i < 7 {
			if x == 0 {
				maxr = (minr + maxr) / 2
			} else {
				minr = (minr + maxr) / 2
			}
		} else {
			if x == 0 {
				maxc = (minc + maxc) / 2
			} else {
				minc = (minc + maxc) / 2
			}
		}
	}

	return minr, minc
}

func getOps(val string) []int {
	vals := []int{}
	for _, s := range val {
		dir := 0
		switch s {
		case 'F', 'L':
			dir = 0
		case 'B', 'R':
			dir = 1
		}
		vals = append(vals, dir)
	}
	return vals
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
