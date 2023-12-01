package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xyz struct {
	x, y, z int
}

type block struct {
	pt          xyz
	left, right *block
	top, bottom *block
	front, back *block

	leftint, rightint bool
	topint, bottomint bool
	frontint, backint bool
}

func (b *block) String() string {
	s := fmt.Sprintf("(%d,%d,%d)\n", b.pt.x, b.pt.y, b.pt.z)
	s += fmt.Sprintf(" left: %t right: %t\n", b.leftint, b.rightint)
	s += fmt.Sprintf(" top: %t bottom: %t\n", b.topint, b.bottomint)
	s += fmt.Sprintf(" front: %t back: %t\n", b.frontint, b.backint)
	return s
}

func print(graph map[xyz]*block) {
	_, max := minmax(graph)
	for z := 0; z <= max.z; z++ {
		fmt.Println("z=", z, "max", max)
		for y := 0; y <= max.y; y++ {
			for x := 0; x <= max.x; x++ {
				if _, ok := graph[xyz{x, y, z}]; ok {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
	}
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	pts := parseInput(in)
	graph := getGraph(pts)
	return getExposedArea(graph)
}

func part2(in input) output {
	pts := parseInput(in)
	graph := getGraph(pts)
	determineExternal(graph)
	return getExposedArea(graph)
}

func graphAny(graph map[xyz]*block, test func(pt xyz, b *block) bool) bool {
	for pt, lava := range graph {
		if test(pt, lava) {
			return true
		}
	}
	return false
}

func getExposedArea(graph map[xyz]*block) int {
	total := 0
	for _, lava := range graph {
		if lava.top == nil && !lava.topint {
			total++
		}
		if lava.bottom == nil && !lava.bottomint {
			total++
		}
		if lava.left == nil && !lava.leftint {
			total++
		}
		if lava.right == nil && !lava.rightint {
			total++
		}
		if lava.front == nil && !lava.frontint {
			total++
		}
		if lava.back == nil && !lava.backint {
			total++
		}
	}
	return total
}

func determineExternal(graph map[xyz]*block) {
	min, max := minmax(graph)
	find := make(map[xyz]bool)
	for z := min.z - 1; z < max.z+1; z++ {
		for y := min.y - 1; y < max.y+1; y++ {
			for x := min.x - 1; x < max.x+1; x++ {
				find[xyz{min.x, y, z}] = true
				find[xyz{x, min.y, z}] = true
				find[xyz{x, y, min.z}] = true
				find[xyz{max.x, y, z}] = true
				find[xyz{x, max.y, z}] = true
				find[xyz{x, y, max.z}] = true
			}
		}
	}
	for bpt, lava := range graph {
		if lava.left == nil && !hasPath(xyz{bpt.x - 1, bpt.y, bpt.z}, find, graph) {
			lava.leftint = true
		}

		if lava.right == nil && !hasPath(xyz{bpt.x + 1, bpt.y, bpt.z}, find, graph) {
			lava.rightint = true
		}

		if lava.top == nil && !hasPath(xyz{bpt.x, bpt.y + 1, bpt.z}, find, graph) {
			lava.topint = true
		}

		if lava.bottom == nil && !hasPath(xyz{bpt.x, bpt.y - 1, bpt.z}, find, graph) {
			lava.bottomint = true
		}

		if lava.front == nil && !hasPath(xyz{bpt.x, bpt.y, bpt.z - 1}, find, graph) {
			lava.frontint = true
		}

		if lava.back == nil && !hasPath(xyz{bpt.x, bpt.y, bpt.z + 1}, find, graph) {
			lava.backint = true
		}
	}
}

func minmax(graph map[xyz]*block) (xyz, xyz) {
	var min, max xyz = xyz{math.MaxInt32, math.MaxInt32, math.MaxInt32}, xyz{math.MinInt32, math.MinInt32, math.MinInt32}
	for pt := range graph {
		if pt.x < min.x {
			min.x = pt.x
		}
		if pt.y < min.y {
			min.y = pt.y
		}
		if pt.z < min.z {
			min.z = pt.z
		}
		if pt.x > max.x {
			max.x = pt.x
		}
		if pt.y > max.y {
			max.y = pt.y
		}
		if pt.z > max.z {
			max.z = pt.z
		}
	}
	return min, max
}

func hasPath(from xyz, to map[xyz]bool, graph map[xyz]*block) bool {
	queue := []xyz{from}
	visit := make(map[xyz]bool)
	result := false

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if _, ok := to[cur]; ok {
			result = true
			break
		}

		if _, ok := visit[cur]; ok {
			continue
		}
		visit[cur] = true

		mvs := []xyz{
			{cur.x + 1, cur.y, cur.z},
			{cur.x - 1, cur.y, cur.z},
			{cur.x, cur.y - 1, cur.z},
			{cur.x, cur.y + 1, cur.z},
			{cur.x, cur.y, cur.z - 1},
			{cur.x, cur.y, cur.z + 1},
		}

		for _, mv := range mvs {
			if _, ok := graph[mv]; ok {
				continue
			}
			queue = append(queue, mv)
		}
	}
	return result
}

func getGraph(pts []xyz) map[xyz]*block {
	blocks := make(map[xyz]*block)
	for _, pt := range pts {
		blocks[pt] = &block{pt: pt}
	}

	for pt, b := range blocks {
		top := xyz{pt.x, pt.y + 1, pt.z}
		if _, ok := blocks[top]; ok {
			blocks[top].bottom = b
			blocks[pt].top = blocks[top]
		}

		bottom := xyz{pt.x, pt.y - 1, pt.z}
		if _, ok := blocks[bottom]; ok {
			blocks[bottom].top = b
			blocks[pt].bottom = blocks[bottom]
		}

		left := xyz{pt.x - 1, pt.y, pt.z}
		if _, ok := blocks[left]; ok {
			blocks[left].right = b
			blocks[pt].left = blocks[left]
		}

		right := xyz{pt.x + 1, pt.y, pt.z}
		if _, ok := blocks[right]; ok {
			blocks[right].left = b
			blocks[pt].right = blocks[right]
		}

		front := xyz{pt.x, pt.y, pt.z - 1}
		if _, ok := blocks[front]; ok {
			blocks[front].back = b
			blocks[pt].front = blocks[front]
		}

		back := xyz{pt.x, pt.y, pt.z + 1}
		if _, ok := blocks[back]; ok {
			blocks[back].front = b
			blocks[pt].back = blocks[back]
		}
	}

	return blocks
}

func keys[K comparable, V any](m map[K]V) []K {
	list := []K{}
	for k := range m {
		list = append(list, k)
	}
	return list
}

func parseInput(in input) []xyz {
	pts := []xyz{}
	for _, line := range in {
		sp := strings.Split(line, ",")
		if len(sp) != 3 {
			continue
		}
		var x, y, z int

		x, _ = strconv.Atoi(sp[0])
		y, _ = strconv.Atoi(sp[1])
		z, _ = strconv.Atoi(sp[2])

		pts = append(pts, xyz{x, y, z})
	}
	return pts
}
