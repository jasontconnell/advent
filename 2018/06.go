package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "06.txt"

var ids string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type coord struct {
	id   string
	x, y int
	inf bool
}

func (c coord) String() string {
	return fmt.Sprintf("%s (%d, %d)", c.id, c.x, c.y)
}

type point struct {
	x, y    int
	nearest *coord
	tied    bool
	dist    int
	inf     bool
}

type xy struct {
	x, y int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	coords := []coord{}
	iid := 0

	for scanner.Scan() {
		var txt = scanner.Text()
		ss := strings.Fields(strings.Replace(txt, ",", "", 1))
		x, _ := strconv.Atoi(ss[0])
		y, _ := strconv.Atoi(ss[1])

		id := string(ids[iid])
		coords = append(coords, coord{id: id, x: x, y: y})
		iid++
	}

	minX, maxX, minY, maxY := minMax(coords)
	fmt.Println(minX, maxX, minY, maxY)
	points := populate(minX, maxX, minY, maxY)
	updated := determineInfs(points, minX, maxX, minY, maxY)
	updated = populateDistances(points, coords)

	id, largest := determineLargest(updated)

	fmt.Println(id, largest)

	fmt.Println("Time", time.Since(startTime))
}

func minMax(coords []coord) (minX, maxX, minY, maxY int){
	minX,minY = 10000,10000
	for _, c := range coords {
		if c.x > maxX {
			maxX = c.x
		}

		if c.y > maxY {
			maxY = c.y
		}

		if c.x < minX {
			minX = c.x
		}

		if c.y < minY {
			minY = c.y
		}
	}

	return
}

func determineInfs(points []point, minX, maxX, minY, maxY int) []point {
	for i := 0; i < len(points); i++ {
		p := &points[i]
		p.inf = p.x < minX || p.y < minY || p.x > maxX || p.y > maxY
	}

	return points
}

func determineLargest(points []point) (string, int) {
	m := make(map[string]int)
	for _, p := range points {
		if p.nearest == nil || p.nearest.inf {
			continue
		}

		m[p.nearest.id]++
	}

	maxsize := 0
	var maxid string
	for k, v := range m {
		if v > maxsize {
			maxsize = v
			maxid = k
		}
	}
	return maxid, maxsize
}

func populateDistances(points []point, coords []coord) []point {
	for i := 0; i < len(points); i++ {
		p := &points[i]
		for j := 0; j < len(coords); j++ {
			c := &coords[j]
			d := distance(p.x, p.y, c.x, c.y)
			if d < p.dist {
				p.nearest = c
				p.dist = d
				p.tied = false
			} else if d == p.dist {
				p.nearest = nil
				p.tied = true
			}
		}
		if p.nearest != nil {
			p.nearest.inf = p.inf
		}
	}

	return points
}

func populate(minX, maxX, minY, maxY int) []point {
	points := []point{}
	for i := minX-1; i < maxX+1; i++ {
		for j := minY-1; j < maxY+1; j++ {
			p := point{x: i, y: j, nearest: nil, tied: false, dist: 10000}
			points = append(points, p)
		}
	}
	return points
}

func distance(x, y, x1, y1 int) int {
	dx := math.Abs(float64(x) - float64(x1))
	dy := math.Abs(float64(y) - float64(y1))

	return int(dx + dy)
}
