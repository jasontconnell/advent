package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func (pt xy) add(pt2 xy) xy {
	return xy{pt.x + pt2.x, pt.y + pt2.y}
}

func maxxy(m map[xy]color) (int, int) {
	maxx, maxy := math.MinInt32, math.MinInt32

	for k := range m {
		if k.x > maxx {
			maxx = k.x
		}
		if k.y > maxy {
			maxy = k.y
		}
	}
	return maxx, maxy
}

var (
	left  = xy{-1, 0}
	right = xy{1, 0}
	up    = xy{0, -1}
	down  = xy{0, 1}
)
var dirs = []xy{left, right, up, down}

type color int

const (
	red   color = 0
	green color = 1
)

type area struct {
	p1, p2 xy
	vol    int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 09 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	points := parseInput(in)
	areas := mapAreas(points)
	return areas[len(areas)-1].vol
}

func part2(in input) output {
	points := parseInput(in)
	colored := colorBorder(points)
	fillInner(colored, green)
	log.Println("filled inner")
	return largestColoredArea(points, colored)
}

func mapAreas(pts []xy) []area {
	areas := []area{}
	for i := 0; i < len(pts)-1; i++ {
		for j := i + 1; j < len(pts); j++ {
			p1, p2 := pts[i], pts[j]
			vol := getArea(p1.x, p2.x, p1.y, p2.y)
			areas = append(areas, area{p1, p2, vol})
		}
	}
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].vol < areas[j].vol
	})
	return areas
}

func largestColoredArea(pts []xy, m map[xy]color) int {
	areas := mapAreas(pts)
	vol := 0
	for i := len(areas) - 1; i >= 0; i-- {
		a := areas[i]
		if allColored(m, a.p1, a.p2) {
			vol = a.vol
			break
		}
	}
	return vol
}

func allColored(m map[xy]color, p1, p2 xy) bool {
	cnt := true
	for y := common.Min(p1.y, p2.y); y <= common.Max(p1.y, p2.y) && cnt; y++ {
		for x := common.Min(p1.x, p2.x); x <= common.Max(p1.x, p2.x); x++ {
			if _, ok := m[xy{x, y}]; !ok {
				cnt = false
				break
			}
		}
	}
	return cnt
}

func getArea(x1, x2, y1, y2 int) int {
	return int((math.Abs(float64(x2-x1)) + 1) * (math.Abs(float64(y2-y1)) + 1))
}

func colorBorder(points []xy) map[xy]color {
	m := make(map[xy]color)
	for _, pt := range points {
		m[pt] = red
	}

	for i := 0; i < len(points); i++ {
		cur := points[i]
		pidx := i - 1
		if i == 0 {
			pidx = len(points) - 1
		}
		prev := points[pidx]
		if cur.x != prev.x && cur.y == prev.y { // y is promised to be equal
			for j := common.Min(cur.x, prev.x) + 1; j < common.Max(cur.x, prev.x); j++ {
				pt := xy{j, cur.y}
				m[pt] = green
			}
		} else if cur.y != prev.y && cur.x == prev.x { // same comment but x
			for j := common.Min(cur.y, prev.y) + 1; j < common.Max(cur.y, prev.y); j++ {
				pt := xy{cur.x, j}
				m[pt] = green
			}
		}
	}
	return m
}

func fillInner(m map[xy]color, c color) {
	maxx, maxy := maxxy(m)
	pt := findInner(m)
	def := xy{-1, -1}
	if pt == def {
		log.Println("couldn't find inner")
		return
	}
	v := make(map[xy]bool)
	log.Println("filling from ", pt)

	queue := []xy{pt}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.x < 0 || cur.y < 0 || cur.x > maxx || cur.y > maxy {
			continue
		}

		if _, ok := v[cur]; ok {
			continue
		}
		v[cur] = true

		if _, ok := m[cur]; ok {
			// color if not already a border
			continue
		}

		m[cur] = c

		for _, d := range dirs {
			np := cur.add(d)
			queue = append(queue, np)
		}
	}
}

func findInner(m map[xy]color) xy {
	for k := range m {
		for _, d := range dirs {
			if isInner(m, k.add(d)) {
				return k.add(d)
			}
		}
	}
	return xy{-1, -1}
}

func isInner(m map[xy]color, pt xy) bool {
	maxx, maxy := maxxy(m)
	for _, d := range dirs {
		cp := pt
		for {
			cp = cp.add(d)
			if _, ok := m[cp]; ok {
				break
			}
			if cp.x < 0 || cp.x > maxx {
				return false
			}
			if cp.y < 0 || cp.y > maxy {
				return false
			}
		}
	}

	return true
}

func parseInput(in input) []xy {
	list := []xy{}
	for _, line := range in {
		sp := strings.Split(line, ",")
		x, _ := strconv.Atoi(sp[0])
		y, _ := strconv.Atoi(sp[1])
		list = append(list, xy{x, y})
	}
	return list
}
