package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var input = "03.txt"

type direction int

const (
	left direction = iota
	right
	up
	down
	none
)

type move struct {
	dir   direction
	units int
}

type wirepath struct {
	start point
	moves []move
}

type point struct {
	x, y int
}

type segmentdir int

const (
	horizontal segmentdir = iota
	vertical
)

type segment struct {
	p1, p2 point
	dir    segmentdir
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wirepaths := []wirepath{}
	start := point{0, 0}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var txt = scanner.Text()

		wp := wirepath{start: start}
		md := strings.Split(txt, ",")
		for _, mv := range md {
			wp.moves = append(wp.moves, getMove(mv))
		}
		wirepaths = append(wirepaths, wp)
	}

	if len(wirepaths) != 2 {
		panic("need two wires")
	}

	corners1 := traverse(wirepaths[0])
	corners2 := traverse(wirepaths[1])

	intersections := intersect(corners1, corners2)
	if intersections == nil || len(intersections) == 0 {
		panic("NO POINTS INTERSECT")
	}

	sort.Slice(intersections, func(i, j int) bool {
		return dist(start, intersections[i]) < dist(start, intersections[j])
	})

	fmt.Println("Part 1: ", dist(start, intersections[0])) // first is closest

	fmt.Println("Time", time.Since(startTime))
}

func dist(start, p point) int {
	return int(math.Abs(float64((p.x - start.x) + (p.y - start.y))))
}

func intersect(path1, path2 []segment) []point {
	points := []point{}
	for _, p := range path1 {
		intersects := getIntersects(p, path2)
		points = append(points, intersects...)
	}
	return points
}

func getIntersects(s segment, test []segment) []point {
	intersects := []point{}
	for _, p := range test {
		if s.dir == p.dir && (s.p1.x != s.p2.x || s.p1.y != s.p2.y) { // can't intersect if they are going the same way on different planes
			continue
		}

		var check1, check2 segment

		if s.dir == vertical {
			check1 = s
			check2 = p
		} else {
			check1 = p
			check2 = s
		}

		if (check1.p1.x > check2.p1.x && check1.p1.x > check2.p2.x) ||
			(check1.p1.x < check2.p1.x && check1.p1.x < check2.p2.x) {
			continue
		}

		if (check2.p1.y > check1.p1.y && check2.p1.y > check1.p2.y) ||
			(check2.p1.y < check1.p1.y && check2.p1.y < check1.p2.y) {
			continue
		}

		intersect := point{x: check1.p1.x, y: check2.p1.y}
		intersects = append(intersects, intersect)
	}

	return intersects
}

func traverse(wp wirepath) []segment {
	segments := []segment{}
	last := wp.start
	for _, m := range wp.moves {
		var n point = last
		s := segment{p1: last}
		switch m.dir {
		case left:
			n.x -= m.units
			s.dir = horizontal
		case right:
			n.x += m.units
			s.dir = horizontal
		case up:
			n.y += m.units
			s.dir = vertical
		case down:
			n.y -= m.units
			s.dir = vertical
		}

		s.p2 = n

		segments = append(segments, s)
		last = n
	}

	return segments
}

func getMove(p string) move {
	m := move{}
	m.dir = getDir(p[0])
	u, err := strconv.Atoi(p[1:])
	if err != nil {
		fmt.Println(err)
	}

	m.units = u

	return m
}

func getDir(p byte) direction {
	var d direction
	switch p {
	case 'L':
		d = left
	case 'R':
		d = right
	case 'U':
		d = up
	case 'D':
		d = down
	}

	return d
}
