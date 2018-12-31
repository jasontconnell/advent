package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"math"
	"strconv"
	"strings"
)

var input = "25.txt"

type xyza struct {
	x, y, z, a int
}

type constellation struct {
	points []xyza
}

func (pt xyza) String() string {
	return fmt.Sprintf("(%d, %d, %d, %d)", pt.x, pt.y, pt.z, pt.a)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	pts := []xyza{}
	for scanner.Scan() {
		var txt = scanner.Text()
		pt := getXYZA(txt)
		if pt != nil {
			pts = append(pts, *pt)
		}
	}

	p1 := solve(pts)

	fmt.Println("Part 1:", p1)

	fmt.Println("Time", time.Since(startTime))
}

func solve(points []xyza) int {
	clone := append([]xyza{}, points...)
	count := 0
	visited := make(map[xyza]bool)
	for len(clone) > 0 {
		list := getConstellation(clone[0], clone, visited)
		if len(list) == 0 {
			break
		}
		clone = removePoints(clone, list)
		count++
	}

	return count
}

func removePoints(src, list []xyza) []xyza {
	cp := append([]xyza{}, src...)
	for i := len(cp) - 1; i >= 0; i-- {
		for j := 0; j < len(list); j++ {
			if cp[i] == list[j] {
				cp = append(cp[:i], cp[i+1:]...)
				break
			}
		}
	}
	return cp
}

func getConnections(point xyza, points []xyza) []xyza {
	conns := []xyza{}
	for _, p := range points {
		if distance(point, p) <= 3 {
			conns = append(conns, p)
		}
	}
	return conns
}

func getConstellation(pt xyza, points []xyza, m map[xyza]bool) []xyza {
	conns := getConnections(pt, points)
	if len(conns) == 0 {
		return []xyza{pt}
	}
	list := []xyza{}
	for _, p := range conns {
		list = append(list, p)
		if _, ok := m[p]; !ok {
			m[p] = true
			nlist := getConstellation(p, points, m)
			list = append(list, nlist...)
		}
	}
	return list
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func distance(a, b xyza) int {
	dx := abs(a.x - b.x)
	dy := abs(a.y - b.y)
	dz := abs(a.z - b.z)
	da := abs(a.a - b.a)

	return dx + dy + dz + da
}

func getXYZA(line string) *xyza {
	sp := strings.Split(line, ",")
	if len(sp) == 4 {
		x, _ := strconv.Atoi(sp[0])
		y, _ := strconv.Atoi(sp[1])
		z, _ := strconv.Atoi(sp[2])
		a, _ := strconv.Atoi(sp[3])

		return &xyza{x: x, y: y, z: z, a: a}
	}
	return nil
}
