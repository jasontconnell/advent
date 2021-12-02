package main

import (
	"fmt"
	"os"
)

type vector2d struct {
	x int
	y int
}

type vector3d struct {
	x int
	y int
	z int
}

type state struct {
	pos   vector3d
	steps int
}

func getGrid() (grid map[vector2d]bool, doors, keys map[vector2d]int, start vector2d, allkeys int) {
	input, _ := os.Open("18.txt")
	grid, doors, keys = make(map[vector2d]bool), make(map[vector2d]int), make(map[vector2d]int)
	var x, y int
	for {
		var line string
		_, err := fmt.Fscanln(input, &line)
		if err != nil {
			break
		}

		for _, c := range line {
			if c != '#' {
				grid[vector2d{x, y}] = true
				if c == '@' {
					start = vector2d{x, y}
				} else if c != '.' {
					if c < 'a' {
						k := 1 << (c - 'A')
						doors[vector2d{x, y}] = k
						allkeys |= k
					} else {
						k := 1 << (c - 'a')
						keys[vector2d{x, y}] = k
						allkeys |= k
					}
				}
			}

			x++
		}

		x = 0
		y++
	}

	return
}

func search(grid map[vector2d]bool, doors, keys map[vector2d]int, start vector2d, allkeys, haveKeys int) int {
	directions := []vector2d{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	queue, visited := []state{state{pos: vector3d{start.x, start.y, haveKeys}}}, make(map[vector3d]bool)
	var st state
	for {
		st, queue = queue[0], queue[1:]

		if st.pos.z&allkeys == allkeys {
			return st.steps
		}

		visited[st.pos] = true

		for _, d := range directions {
			next := vector3d{st.pos.x + d.x, st.pos.y + d.y, st.pos.z}

			if !grid[vector2d{next.x, next.y}] || visited[next] {
				continue
			}

			door, ok := doors[vector2d{next.x, next.y}]
			if ok && next.z&door != door {
				continue
			}

			key, ok := keys[vector2d{next.x, next.y}]
			if ok {
				next.z |= key
			}

			queue = append(queue, state{pos: next, steps: st.steps + 1})
		}
	}
}

func part1() {
	grid, doors, keys, start, allkeys := getGrid()
	fmt.Println(search(grid, doors, keys, start, allkeys, 0))
}

func part2() {
	grid, doors, keys, start, allkeys := getGrid()

	directions := []vector2d{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	grid[start] = false
	for _, d := range directions {
		grid[vector2d{start.x + d.x, start.y + d.y}] = false
	}

	total := 0
	haveKeys := allkeys
	for x := 0; x < start.x; x++ {
		for y := 0; y < start.y; y++ {
			haveKeys ^= keys[vector2d{x, y}]
		}
	}

	total += search(grid, doors, keys, vector2d{start.x - 1, start.y - 1}, allkeys, haveKeys)

	haveKeys = allkeys
	for x := start.x + 1; x <= start.x*2; x++ {
		for y := 0; y < start.y; y++ {
			haveKeys ^= keys[vector2d{x, y}]
		}
	}

	total += search(grid, doors, keys, vector2d{start.x + 1, start.y - 1}, allkeys, haveKeys)

	haveKeys = allkeys
	for x := start.x + 1; x <= start.x*2; x++ {
		for y := start.y + 1; y <= start.y*2; y++ {
			haveKeys ^= keys[vector2d{x, y}]
		}
	}

	total += search(grid, doors, keys, vector2d{start.x + 1, start.y + 1}, allkeys, haveKeys)

	haveKeys = allkeys
	for x := 0; x < start.x; x++ {
		for y := start.y + 1; y <= start.y*2; y++ {
			haveKeys ^= keys[vector2d{x, y}]
		}
	}

	total += search(grid, doors, keys, vector2d{start.x - 1, start.y + 1}, allkeys, haveKeys)

	fmt.Println(total)
}

func main() {
	part := 0
	if len(os.Args) == 2 {
		fmt.Sscan(os.Args[1], &part)
	} else {
		fmt.Print("Enter 1 or 2 to select part: ")
		fmt.Scanf("%d\n", &part)
	}

	switch part {
	case 1:
		part1()
	case 2:
		part2()
	default:
		fmt.Println("Error: Invalid part.")
	}
}
