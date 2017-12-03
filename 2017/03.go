package main

import (
	"fmt"
	"math"
	"time"
)

var input = 368078

func main() {
	startTime := time.Now()

	d := dist(input)
	fmt.Println("distance", d)

	first := fillGrid(input)
	fmt.Println("first filled larger than input", first)

	fmt.Println("Time", time.Since(startTime))
}

func ring(i int) int {
	if i == 1 {
		return 1
	}
	s := 1
	itr := 1

	for (s * s) < i {
		s += 2
		itr++
	}

	return s - 2
}

func coords(i int) (int, int, int) {
	sq := ring(i)
	start := sq*sq + 1
	max := sq + 1
	startx, starty := max, max-1

	diff := i - start
	turns := 0

	extra := diff
	if diff > max-1 {
		for j := 0; j < 4; j++ {
			m := max
			if j == 0 {
				m = m - 1
			}

			if extra-m >= 0 {
				extra = extra - m
				turns++
			}
		}
	}

	x, y := 0, 0
	switch turns {
	case 0:
		x = max
		y = starty - extra
	case 1:
		x = startx - extra
		y = 0
	case 2:
		x = 0
		y = extra
	case 3:
		x = extra
		y = max
	}

	return x, y, max
}

func dist(i int) float64 {
	x, y, max := coords(i)
	one := max / 2 // position of 1

	return math.Abs(float64(x-one)) + math.Abs(float64(y-one))
}

func fillGrid(d int) int {
	val := 0
	size := 1000
	m := [][]int{}

	startx, starty := size/2, size/2

	m = make([][]int, size)

	for i := 0; i < size; i++ {
		m[i] = make([]int, size)
	}

	xdir, ydir := 1, 0

	curx := startx
	cury := starty

	m[cury][curx] = 1

	//
	for i := 0; i < size*size && val == 0; i++ {
		curx = curx + xdir
		cury = cury + ydir

		dur, dul, ddr, ddl := m[cury-1][curx+1], m[cury-1][curx-1], m[cury+1][curx+1], m[cury+1][curx-1]
		up, down, left, right := m[cury-1][curx], m[cury+1][curx], m[cury][curx-1], m[cury][curx+1]
		m[cury][curx] = dur + dul + ddr + ddl + up + down + left + right

		if m[cury][curx] > d {
			val = m[cury][curx]
			break
		}

		// turn??
		if xdir == 1 { // going right. if up is 0, turn up
			if up == 0 {
				xdir = 0
				ydir = -1
			}
		} else if xdir == -1 { // going left. if down is 0, turn down
			if down == 0 {
				xdir = 0
				ydir = 1
			}
		} else if ydir == 1 { // going down. if right is 0, turn right
			if right == 0 {
				xdir = 1
				ydir = 0
			}
		} else if ydir == -1 { // going up. if left is 0, turn left
			if left == 0 {
				xdir = -1
				ydir = 0
			}
		}
	}

	return val
}
