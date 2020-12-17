package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "17.txt"

type point3d struct {
	x, y, z int
}

type point4d struct {
	x, y, z, w int
}

func (p point3d) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.x, p.y, p.z)
}

var matrix []point3d
var matrix4d []point4d

func init() {
	matrix = permMatrix()
	matrix4d = permMatrix4d()
}

func permMatrix() []point3d {
	list := []point3d{}
	visited := make(map[point3d]bool)
	// center piece isn't a neighbor
	visited[point3d{0, 0, 0}] = true

	for z := -1; z <= 1; z++ {
		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				p := point3d{x, y, z}

				if _, ok := visited[p]; !ok {
					list = append(list, p)
					visited[p] = true
				}
			}
		}
	}

	return list
}

func permMatrix4d() []point4d {
	list := []point4d{}
	visited := make(map[point4d]bool)
	// center piece isn't a neighbor
	visited[point4d{0, 0, 0, 0}] = true

	for w := -1; w <= 1; w++ {
		for z := -1; z <= 1; z++ {
			for y := -1; y <= 1; y++ {
				for x := -1; x <= 1; x++ {
					p := point4d{x, y, z, w}
					if _, ok := visited[p]; !ok {
						list = append(list, p)
						visited[p] = true
					}
				}
			}
		}
	}

	return list
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

	m := readMap(lines)
	mc := simulate3d(m, 6)

	p1 := countOn3d(mc)

	fmt.Println("Part 1:", p1)

	m2 := readMap(lines)
	m4d := make(map[point4d]bool)
	for k, v := range m2 {
		p4d := point4d{x: k.x, y: k.y, z: 0, w: 0}
		m4d[p4d] = v
	}

	mp2 := simulate4d(m4d, 6)
	p2 := countOn4d(mp2)

	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func countOn3d(m map[point3d]bool) int {
	count := 0
	for _, v := range m {
		if v {
			count++
		}
	}
	return count
}

func countOn4d(m map[point4d]bool) int {
	count := 0
	for _, v := range m {
		if v {
			count++
		}
	}
	return count
}

func simulate3d(m map[point3d]bool, count int) map[point3d]bool {
	mc := make(map[point3d]bool)
	for k, v := range m {
		mc[k] = v
	}

	for i := 0; i < count; i++ {
		mc = simulateOne3d(mc)
	}
	return mc
}

func simulate4d(m map[point4d]bool, count int) map[point4d]bool {
	mc := make(map[point4d]bool)
	for k, v := range m {
		mc[k] = v
	}

	for i := 0; i < count; i++ {
		mc = simulateOne4d(mc)
	}
	return mc
}

func simulateOne3d(m map[point3d]bool) map[point3d]bool {
	mc := make(map[point3d]bool, len(m))

	for k, v := range m {
		if v {
			for _, n := range getNeighbors3d(k) {
				if _, ok := m[n]; !ok {
					m[n] = false
				}
			}
		}
	}

	// read over m, write to mc
	for k, active := range m {
		oncount := onCount3d(m, k)

		switch oncount {
		case 2, 3:
			if !active && oncount == 3 {
				mc[k] = true
			} else if active {
				mc[k] = true
			}
		default:
			if active {
				mc[k] = false
			}
		}
	}

	return mc
}

func simulateOne4d(m map[point4d]bool) map[point4d]bool {
	mc := make(map[point4d]bool, len(m))

	for k, v := range m {
		if v {
			for _, n := range getNeighbors4d(k) {
				if _, ok := m[n]; !ok {
					m[n] = false
				}
			}
		}
	}

	// read over m, write to mc
	for k, active := range m {
		oncount := onCount4d(m, k)

		switch oncount {
		case 2, 3:
			if !active && oncount == 3 {
				mc[k] = true
			} else if active {
				mc[k] = true
			}
		default:
			if active {
				mc[k] = false
			}
		}

		// if mc[k] {
		// 	for _, n := range getNeighbors4d(k) {
		// 		if _, ok := mc[n]; !ok {
		// 			mc[n] = false
		// 		}
		// 	}
		// }
	}

	return mc
}

func getNeighbors3d(p point3d) []point3d {
	list := []point3d{}
	for _, mp := range matrix {
		np := point3d{x: p.x + mp.x, y: p.y + mp.y, z: p.z + mp.z}
		list = append(list, np)
	}
	return list
}

func getNeighbors4d(p point4d) []point4d {
	list := []point4d{}
	for _, mp := range matrix4d {
		np := point4d{x: p.x + mp.x, y: p.y + mp.y, z: p.z + mp.z, w: p.w + mp.w}
		list = append(list, np)
	}
	return list
}

func onCount3d(m map[point3d]bool, p point3d) int {
	on := 0
	for _, n := range getNeighbors3d(p) {
		if isOn, ok := m[n]; ok && isOn {
			on++
		}
	}
	return on
}

func onCount4d(m map[point4d]bool, p point4d) int {
	on := 0
	for _, n := range getNeighbors4d(p) {
		if isOn, ok := m[n]; ok && isOn {
			on++
		}
	}
	return on
}

func readMap(lines []string) map[point3d]bool {
	m := make(map[point3d]bool)
	for y, line := range lines {
		for x, ch := range line {
			p := point3d{x: x, y: y, z: 0}
			on := ch == '#'

			m[p] = on
		}
	}

	return m
}
